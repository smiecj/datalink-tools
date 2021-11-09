package main

import (
	"fmt"
	"strings"
	"sync"

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
	clientSingleton Client
	clientOnce      sync.Once
)

type Client interface {
	GetMedias() ([]*Media, error)
	GetTasks() ([]*Task, error)
	GetMappings() ([]*Mapping, error)
	login() error
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
	log.Info("[login] 登录 Azkaban 结果: %s", rsp.Body)
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

// 获取所有介质
// todo: 当前: 调试参数错误的问题: ① length 不能为0 ② ？
func (client *datalinkClient) GetMedias() ([]*Media, error) {
	retMediaArr := make([]*Media, 0)
	// 分页获取
	for {
		start := 0
		queryParam := buildQueryParamWithPage(start, pageStep)
		queryParamStr, _ := json.Marshal(queryParam)
		getMediaUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlGetMedia)
		rsp, err := client.httpClient.Do(http.Url(getMediaUrl),
			http.Post(),
			http.AddHeader(reqHeaderCookie, fmt.Sprintf(reqCookieSessionIdFormat, client.sessionId)),
			http.SetParam(string(queryParamStr)))
		if nil != err {
			log.Error("[GetMedias] 获取介质失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		queryMediaRet := QueryMediaListRet{}
		err = json.Unmarshal([]byte(rsp.Body), &queryMediaRet)
		if nil != err {
			log.Error("[GetMedias] 解析获取介质结果失败, 当前起点位置: %d，错误信息: %s", start, err.Error())
			return retMediaArr, err
		}

		if len(queryMediaRet.MediaList) == 0 {
			log.Info("[GetMedias] 获取介质完成, 起点位置: %d", start)
			return retMediaArr, nil
		}

		for _, currentMedia := range queryMediaRet.MediaList {
			retMediaArr = append(retMediaArr, &currentMedia)
		}

		start += pageStep
	}
}

// 获取所有任务
func (client *datalinkClient) GetTasks() ([]*Task, error) {
	// todo: implement
	return nil, fmt.Errorf("not implement")
}

// 获取所有同步关联配置
func (client *datalinkClient) GetMappings() ([]*Mapping, error) {
	// todo: implement
	return nil, fmt.Errorf("not implement")
}

// 获取 client 单例
func GetDataLinkClient(option DatalinkOption) Client {
	clientOnce.Do(func() {
		client := new(datalinkClient)
		client.Option = option
		client.httpClient = http.GetHTTPClient()
		clientSingleton = client
	})
	return clientSingleton
}
