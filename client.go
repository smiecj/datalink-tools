package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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
)

var (
	clientMap  = make(map[DatalinkOption]Client)
	clientLock sync.RWMutex
)

type Client interface {
	GetRDBMedias() ([]*RDBMedia, error)
	GetKuduMedias() ([]*KuduMedia, error)
	GetTasks() ([]*Task, error)
	GetMappings() ([]*Mapping, error)
	SaveMedias([]*RDBMedia) (int, error)
	SaveTasks([]*Task) (int, error)
	SaveMappings([]*Mapping) (int, error)
	start()
}

type datalinkClient struct {
	Option     DatalinkOption
	sessionId  string
	httpClient http.Client
}

// 登录
func (client *datalinkClient) login() error {
	// 调用 login 接口
	loginUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlLogin)
	rsp, err := client.httpClient.Do(http.Url(loginUrl),
		http.PostWithUrlEncode(),
		http.AddParam("loginEmail", client.Option.Username),
		http.AddParam("password", client.Option.Password))
	log.Info("[login] 登录 datalink 结果: %s", rsp.Body)
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
		getRDBMediaUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlGetRDBMedia)
		rsp, err := client.httpClient.Do(http.Url(getRDBMediaUrl),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			log.Error("[GetRDBMedias] 获取RDB介质失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		queryMediaRet := QueryRDBMediaListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryMediaRet)
		if nil != err {
			log.Error("[GetRDBMedias] 解析获取RDB介质结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		if len(queryMediaRet.MediaList) == 0 {
			log.Info("[GetRDBMedias] 获取RDB介质完成, 起点位置: %d", start)
			return retMediaArr, nil
		}

		for _, currentMedia := range queryMediaRet.MediaList {
			retMediaArr = append(retMediaArr, &currentMedia)
		}
		log.Info("[GetRDBMedias] 获取RDB介质完成，起点: %d, 获取介质数: %d", start, len(queryMediaRet.MediaList))

		start += pageStep
	}
}

// 获取所有Kudu 介质
func (client *datalinkClient) GetKuduMedias() ([]*KuduMedia, error) {
	retArr := make([]*KuduMedia, 0)
	// kudu 介质: 直接获取
	getKuduMediaUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlGetKuduMedia)
	rsp, err := client.httpClient.Do(http.Url(getKuduMediaUrl),
		http.Get(),
		http.AddHeader(reqHeaderCookie, client.getCookieSessionId()))
	if nil != err {
		log.Error("[GetKuduMedias] 获取 kudu 介质失败，请检查: %s", err.Error())
		return retArr, err
	}
	queryKuduMediaRet := QueryKuduMediaListRet{}
	err = json.Unmarshal([]byte(rsp.Body), &queryKuduMediaRet)
	if nil != err {
		log.Error("[GetKuduMedias] 解析获取kudu介质结果失败，错误信息: %s", err.Error())
		return retArr, err
	}

	for _, currentKuduMedia := range queryKuduMediaRet.MediaList {
		retArr = append(retArr, &currentKuduMedia)
	}

	log.Info("[GetKuduMedias] 获取 kudu 介质成功，总数: %d", len(retArr))
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
		getTasksUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlGetTasks)
		rsp, err := client.httpClient.Do(http.Url(getTasksUrl),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			log.Error("[GetTasks] 获取任务失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retTasksArr, err
		}

		queryTaskRet := QueryTaskListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryTaskRet)
		if nil != err {
			log.Error("[GetTasks] 解析获取任务结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retTasksArr, err
		}

		if len(queryTaskRet.TaskList) == 0 {
			log.Info("[GetTasks] 获取任务完成, 当前位置: %d", start)
			return retTasksArr, nil
		}

		for _, currentTask := range queryTaskRet.TaskList {
			retTasksArr = append(retTasksArr, &currentTask)
		}
		log.Info("[GetMedias] 获取介质完成，起点: %d, 获取介质数: %d", start, len(queryTaskRet.TaskList))

		start += pageStep
	}
}

// 获取所有同步关联配置
func (client *datalinkClient) GetMappings() ([]*Mapping, error) {
	retMappingArr := make([]*Mapping, 0)
	// 分页获取
	start := 0
	for {
		queryParam := client.buildQueryMappingParam(start, pageStep)
		queryParamStr, _ := json.Marshal(queryParam)
		getMappingUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlGetMapping)
		rsp, err := client.httpClient.Do(http.Url(getMappingUrl),
			http.Post(),
			http.AddHeader(reqHeaderCookie, client.getCookieSessionId()),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			log.Error("[GetMappings] 获取映射失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMappingArr, err
		}

		queryMappingRet := QueryMappingListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryMappingRet)
		if nil != err {
			log.Error("[GetMappings] 解析获取映射结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMappingArr, err
		}

		if len(queryMappingRet.MappingList) == 0 {
			log.Info("[GetMappings] 获取映射完成, 当前位置: %d", start)
			return retMappingArr, nil
		}

		for _, currentMapping := range queryMappingRet.MappingList {
			retMappingArr = append(retMappingArr, &currentMapping)
		}
		log.Info("[GetMappings] 获取映射完成，起点: %d, 获取介质数: %d", start, len(queryMappingRet.MappingList))

		start += pageStep
	}
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

// 开启一个协程，定期更新session id
func (client *datalinkClient) start() {
	// 登录 属于基本功能，一般不做返回值检查
	_ = client.login()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			_ = client.login()
		}
	}()
}

// 根据配置 获取 client 单例
func GetDataLinkClient(option DatalinkOption) Client {
	var client Client
	clientLock.RLock()
	client = clientMap[option]
	clientLock.RUnlock()

	if nil != client {
		return client
	}

	clientLock.Lock()
	defer clientLock.Unlock()

	dlinkClient := new(datalinkClient)
	dlinkClient.Option = option
	dlinkClient.httpClient = http.GetHTTPClient()
	dlinkClient.start()
	clientMap[option] = dlinkClient

	return dlinkClient
}
