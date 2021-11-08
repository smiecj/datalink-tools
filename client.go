package main

import (
	"fmt"
	"strings"
	"sync"

	http "github.com/smiecj/go_common/http"
	"github.com/smiecj/go_common/util/log"
)

const (
	loginSuccessRet          = "\"success\""
	loginRspHeaderCookie     = "Set-Cookie"
	loginRspCookieJsessionId = "JSESSIONID"
)

var (
	clientSingleton Client
	clientOnce      sync.Once
)

type Client interface {
	GetMedias() ([]*Media, error)
	GetTasks() ([]*Task, error)       // todo: task 定义
	GetMappings() ([]*Mapping, error) // todo: 映射定义
	login() error                     // todo: 实现登录接口
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
			if strings.Contains(currentCookie, loginRspCookieJsessionId) {
				jsessionIdSplitArr := strings.Split(currentCookie, loginRspCookieJsessionId)
				// 这里一般不需要判空，只要登录成功，框架一定会返回 session id
				client.sessionId = jsessionIdSplitArr[1]
				return nil
			}
		}
	}
	return fmt.Errorf("login failed")
}

// 获取所有介质
func (client *datalinkClient) GetMedias() ([]*Media, error) {
	// todo: implement
	return nil, fmt.Errorf("not implement")
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
