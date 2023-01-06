package tools

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/errorcode"
	http "github.com/smiecj/go_common/http"
	"github.com/smiecj/go_common/util/json"
	"github.com/smiecj/go_common/util/log"
)

const (
	loginSuccessRet          = "\"success\""
	loginRspHeaderCookie     = "Set-Cookie"
	loginRspCookieSessionId  = "JSESSIONID"
	reqHeaderCookie          = "Cookie"
	reqCookieSessionIdFormat = "JSESSIONID=%s"
	pageStep                 = 10

	paramLoginEmail = "loginEmail"
	paramPassword   = "password"
	paramId         = "id"
)

var (
	clientMap  = make(map[datalinkConfig]Client)
	clientLock sync.RWMutex
)

type Client interface {
	GetRDBMedias() ([]*RDBMedia, error)
	GetKuduMedias() ([]*KuduMedia, error)
	GetTasks() ([]*Task, error)
	GetTask(int) (*TaskDetail, error)
	GetMappings() ([]*Mapping, error)
	SaveMedias([]*RDBMedia) (int, error)
	SaveTasks([]*Task) (int, error)
	SaveMappings([]*Mapping) (int, error)
	StartTask(int) error
	StopTask(int) error
	Refresh() error
	UpdateTask(*TaskDetail) error
	start()
}

type datalinkClient struct {
	conf               datalinkConfig
	sessionId          string
	httpClient         http.Client
	taskNodeAddressArr []string
	log                log.Logger
}

// 登录
func (client *datalinkClient) login() error {
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlLogin()),
		http.PostWithUrlEncode(),
		http.AddParam(paramLoginEmail, client.conf.User),
		http.AddParam(paramPassword, client.conf.Password))
	client.log.Info("[login] 登录 datalink 结果: %s", rsp.Body)
	if nil != err {
		return err
	}

	// 解析接口返回的 header 中的 session id
	if loginSuccessRet == rsp.Body {
		cookieStr := rsp.Header[loginRspHeaderCookie]
		cookieSplitArr := strings.Split(cookieStr, ";")
		for _, currentCookie := range cookieSplitArr {
			if strings.Contains(currentCookie, loginRspCookieSessionId) {
				jsessionIdSplitArr := strings.Split(currentCookie, "=")
				// 这里一般不需要判空，只要登录成功，框架一定会返回 session id
				client.sessionId = jsessionIdSplitArr[1]
				return nil
			}
		}
	}
	return fmt.Errorf("login failed")
}

// 加载 task 节点 api 地址
func (client *datalinkClient) loadWorkerNodeAddress() error {
	// worker 数量目前不会太多，暂不需要分页获取
	queryParamStr := client.buildQueryWorkerParamStr(0, 10)
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlGetWorker()),
		http.Post(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
		http.SetParam(queryParamStr))
	if nil != err {
		client.log.Error("[loadTaskNodeAddress] 获取worker信息失败，请检查: %s", err.Error())
		return err
	}
	queryWorkerRet := queryWorkerRet{}
	_ = json.Unmarshal([]byte(rsp.Body), &queryWorkerRet)
	client.taskNodeAddressArr = queryWorkerRet.getWorkerAddressArr()
	return nil
}

// 获取作为cookie 参数中的 session id 格式
func (client *datalinkClient) getCookieSessionId() string {
	return fmt.Sprintf(reqCookieSessionIdFormat, client.sessionId)
}

// 获取所有RDB介质
func (client *datalinkClient) GetRDBMedias() ([]*RDBMedia, error) {
	retMediaArr := make([]*RDBMedia, 0)
	// 分页获取
	start := 0
	for {
		queryParam := client.buildQueryMediaParam(start, pageStep)
		queryParamStr, _ := json.Marshal(queryParam)
		rsp, err := client.httpClient.Do(http.Url(client.conf.urlRDBMedia()),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			client.log.Error("[GetRDBMedias] 获取RDB介质失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		queryMediaRet := QueryRDBMediaListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryMediaRet)
		if nil != err {
			client.log.Error("[GetRDBMedias] 解析获取RDB介质结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		if len(queryMediaRet.MediaList) == 0 {
			client.log.Info("[GetRDBMedias] 获取RDB介质完成, 起点位置: %d", start)
			return retMediaArr, nil
		}

		for index := 0; index < len(queryMediaRet.MediaList); index++ {
			retMediaArr = append(retMediaArr, &queryMediaRet.MediaList[index])
		}
		client.log.Info("[GetRDBMedias] 获取RDB介质完成，起点: %d, 获取介质数: %d", start, len(queryMediaRet.MediaList))

		start += pageStep
	}
}

// 获取所有 Kudu 介质
func (client *datalinkClient) GetKuduMedias() ([]*KuduMedia, error) {
	retArr := make([]*KuduMedia, 0)
	// kudu 介质: 直接获取
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlKuduMedia()),
		http.Get(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	if nil != err {
		client.log.Error("[GetKuduMedias] 获取 kudu 介质失败，请检查: %s", err.Error())
		return retArr, err
	}
	queryKuduMediaRet := QueryKuduMediaListRet{}
	err = json.Unmarshal([]byte(rsp.Body), &queryKuduMediaRet)
	if nil != err {
		client.log.Error("[GetKuduMedias] 解析获取kudu介质结果失败，错误信息: %s", err.Error())
		return retArr, err
	}

	for _, currentKuduMedia := range queryKuduMediaRet.MediaList {
		retArr = append(retArr, &currentKuduMedia)
	}

	client.log.Info("[GetKuduMedias] 获取 kudu 介质成功，总数: %d", len(retArr))
	return retArr, nil
}

// 获取所有任务
func (client *datalinkClient) GetTasks() ([]*Task, error) {
	retTasksArr := make([]*Task, 0)
	// 分页获取
	start := 0
	for {
		queryParam := client.buildQueryTaskParam(start, pageStep)
		queryParamStr, _ := json.Marshal(queryParam)
		rsp, err := client.httpClient.Do(http.Url(client.conf.urlTasks()),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			client.log.Error("[GetTasks] 获取任务失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retTasksArr, err
		}

		queryTaskRet := QueryTaskListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryTaskRet)
		if nil != err {
			client.log.Error("[GetTasks] 解析获取任务结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retTasksArr, err
		}

		if len(queryTaskRet.TaskList) == 0 {
			client.log.Info("[GetTasks] 获取同步任务完成，总数: %d", len(retTasksArr))
			return retTasksArr, nil
		}

		for index := 0; index < len(queryTaskRet.TaskList); index++ {
			retTasksArr = append(retTasksArr, &queryTaskRet.TaskList[index])
		}

		start += pageStep
	}
}

// 获取指定任务
func (client *datalinkClient) GetTask(id int) (*TaskDetail, error) {
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlGetTask(id)),
		http.Get(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	if nil != err {
		client.log.Error("[GetTask] 获取指定任务 id: %d 失败，请检查: %s", id, err.Error())
		return nil, err
	}
	taskDetail := TaskDetail{}
	err = json.Unmarshal([]byte(rsp.Body), &taskDetail)
	if nil != err {
		client.log.Error("[GetTask] 获取指定任务 id: %d 失败，请检查: %s", id, err.Error())
		return nil, err
	}
	client.log.Info("[GetTask] 获取任务 id: %d 详细信息成功", id)
	return &taskDetail, nil
}

// 获取所有同步关联配置
func (client *datalinkClient) GetMappings() ([]*Mapping, error) {
	retMappingArr := make([]*Mapping, 0)
	// 分页获取
	start := 0
	for {
		queryParam := client.buildQueryMappingParam(start, pageStep)
		queryParamStr, _ := json.Marshal(queryParam)
		rsp, err := client.httpClient.Do(http.Url(client.conf.urlMapping()),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			client.log.Error("[GetMappings] 获取映射失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMappingArr, err
		}

		queryMappingRet := QueryMappingListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryMappingRet)
		if nil != err {
			client.log.Error("[GetMappings] 解析获取映射结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMappingArr, err
		}

		if len(queryMappingRet.MappingList) == 0 {
			client.log.Info("[GetMappings] 获取映射完成, 总数: %d", len(retMappingArr))
			return retMappingArr, nil
		}

		for index := 0; index < len(queryMappingRet.MappingList); index++ {
			retMappingArr = append(retMappingArr, &queryMappingRet.MappingList[index])
		}
		client.log.Info("[GetMappings] 获取映射完成，起点: %d, 获取介质数: %d", start, len(queryMappingRet.MappingList))

		start += pageStep
	}
}

// 启动指定任务
func (client *datalinkClient) StartTask(taskId int) error {
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlStartTask(taskId)),
		http.PostWithUrlEncode(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	client.log.Info("[StartTask] 启动任务结果: %s", rsp.Body)
	if nil != err {
		return err
	}
	return nil
}

// 停止指定任务
func (client *datalinkClient) StopTask(taskId int) error {
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlStopTask(taskId)),
		http.PostWithUrlEncode(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	client.log.Info("[StopTask] 停止任务结果: %s", rsp.Body)
	if nil != err {
		return err
	}
	return nil
}

// 更新指定任务
func (client *datalinkClient) UpdateTask(task *TaskDetail) error {
	rsp, err := client.httpClient.Do(http.Url(client.conf.urlUpdateTask()),
		http.Post(),
		http.SetBody(task.getDetailUpdate()),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	client.log.Info("[UpdateTask] 更新任务结果: %s", rsp.Body)
	if nil != err {
		return err
	}
	return nil
}

// 刷新 （kudu 元数据）
func (client *datalinkClient) Refresh() error {
	// 获取 kudu media
	kuduMediaArr, err := client.GetKuduMedias()
	if nil != err {
		return err
	}
	// 默认 kudu media 只有一个
	if len(kuduMediaArr) != 1 {
		return errorcode.BuildErrorWithMsg(errorcode.ServiceError, "refresh failed: kudu media size is not 1")
	}

	// 调用 task 节点的刷新接口
	for _, currentWorkerAddress := range client.taskNodeAddressArr {
		_, err := client.httpClient.Do(http.Url(client.conf.urlRefreshKuduMedia(currentWorkerAddress, kuduMediaArr[0].Id)),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
		if nil != err {
			client.log.Warn("[Refresh] worker 节点: %s, 刷新 kudu 元数据失败: %s", currentWorkerAddress, err.Error())
			return err
		}
	}
	return nil
}

// todo: 保存介质信息到服务器中
func (client *datalinkClient) SaveMedias([]*RDBMedia) (int, error) {
	return 0, fmt.Errorf("Not implement")
}

// todo: 保存任务信息到服务器中
func (client *datalinkClient) SaveTasks([]*Task) (int, error) {
	return 0, fmt.Errorf("Not implement")
}

// todo: 保存映射信息到服务器中
func (client *datalinkClient) SaveMappings([]*Mapping) (int, error) {
	return 0, fmt.Errorf("Not implement")
}

// 构建介质查询参数
func (client *datalinkClient) buildQueryMediaParam(start, limit int) QueryRDBMediaParam {
	queryParamStr := fmt.Sprintf(defaultQueryRDBMediaStrFormat, start, limit)
	queryParam := QueryRDBMediaParam{}
	json.Unmarshal([]byte(queryParamStr), &queryParam)
	return queryParam
}

// 构建任务查询参数
func (client *datalinkClient) buildQueryTaskParam(start, limit int) QueryTaskParam {
	queryParamStr := fmt.Sprintf(defaultQueryTaskStrFormat, start, limit)
	queryParam := QueryTaskParam{}
	json.Unmarshal([]byte(queryParamStr), &queryParam)
	return queryParam
}

// 构建映射查询参数
func (client *datalinkClient) buildQueryMappingParam(start, limit int) QueryMappingParam {
	queryParamStr := fmt.Sprintf(defaultQueryMappingStrFormat, start, limit)
	queryParam := QueryMappingParam{}
	json.Unmarshal([]byte(queryParamStr), &queryParam)
	return queryParam
}

// 构建worker查询参数
func (client *datalinkClient) buildQueryWorkerParamStr(start, limit int) string {
	queryParamStr := fmt.Sprintf(defaultQueryWorkerStrFormat, start, limit)
	return queryParamStr
}

// 构建RDB介质保存参数
// {media_name}、{db_write_host}、{db_read_host}、{db_write_user}、{db_read_user}
// {db_write_password}、{db_read_password}
func (client *datalinkClient) buildSaveRDBMediaParam(conf *RDBMediaBackupConfig) string {
	replacer := strings.NewReplacer("{media_name}", conf.Name, "{db_write_host}", conf.DBWriteHost,
		"{db_read_host}", conf.DBReadHost, "{db_write_user}", conf.DBWriteUser, "{db_read_user}", conf.DBReadUser,
		"{db_write_password}", conf.DBWritePassword, "{db_read_password}", conf.DBReadPassword)
	paramStr := replacer.Replace(defaultSaveRDBMediaStrFormat)
	return paramStr
}

// 构建Kudu介质保存参数
// {kudu_master_host}、{kudu_master_port}、{name}、{desc}、{db_name}
func (client *datalinkClient) buildSaveKuduMediaParam(conf *KuduMediaBackupConfig) string {
	replacer := strings.NewReplacer("{name}", conf.Name, "{kudu_master_host}", conf.KuduMasterHost,
		"{kudu_master_port}", strconv.Itoa(conf.KuduMasterPort), "{db_name}", conf.DBName, "{desc}", conf.Desc)
	paramStr := replacer.Replace(defaultSaveKuduMediaStrFormat)
	return paramStr
}

// 开启一个协程，定期更新 session id
func (client *datalinkClient) start() {
	// 登录 属于基本功能，一般不做返回值检查
	_ = client.login()
	_ = client.loadWorkerNodeAddress()

	go func() {
		ticker := time.NewTicker(time.Duration(client.conf.SessionInternal) * time.Second)
		for range ticker.C {
			_ = client.login()
			_ = client.loadWorkerNodeAddress()
		}
	}()
}

// 根据配置 获取 client 单例
func GetDataLinkClient(configManager config.Manager) Client {
	conf := datalinkConfig{}
	space, _ := configManager.GetSpace(datalinkSpaceName)
	err := space.Unmarshal(&conf)
	if nil != err {
		log.Warn("[GetDataLinkClient] get datalink client fail: " + err.Error())
		return nil
	}

	var client Client
	clientLock.RLock()
	client = clientMap[conf]
	clientLock.RUnlock()

	if nil != client {
		return client
	}

	clientLock.Lock()
	defer clientLock.Unlock()

	dlinkClient := new(datalinkClient)
	dlinkClient.conf = conf
	dlinkClient.httpClient = http.GetHTTPClient()
	dlinkClient.log = log.PrefixLogger("datalink")
	dlinkClient.start()
	clientMap[conf] = dlinkClient

	return dlinkClient
}
