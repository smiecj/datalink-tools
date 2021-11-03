package main

import (
	"fmt"

	"github.com/smiecj/go_common/http"
)

type Client interface {
	GetStorages() []*Media
	GetTasks() []*Task      // todo: task 定义
	GetLink() []*Mapping    // todo: 映射定义
	login() (string, error) // todo: 实现登录接口
}

type datalinkClient struct {
	Option DatalinkOption
}

// 登录
// 返回: jsessionid 、error
func (client *datalinkClient) login() (jsessionId string, err error) {
	// todo: 调用azkaban 接口测试
	loginUrl := fmt.Sprintf("http://%s%s", client.Option.Address, urlLogin)
	// todo: common 实现对参数的封装，不要把 map 开放出来
	param := make(map[string]string)
	param["loginEmail"] = client.Option.Username
	param["password"] = client.Option.Password
	http.DoPostFormRequest(loginUrl, param)

	// todo: 在本地进行 azkaban login 接口调用验证
}

func GetDataLinkClient(option DatalinkOption) *datalinkClient {
	client := new(datalinkClient)
	client.Option = option
	return client
}
