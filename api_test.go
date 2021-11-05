package main

import (
	"testing"

	http "github.com/smiecj/go_common/http"
	"github.com/smiecj/go_common/util/log"
)

func TestLogin(t *testing.T) {
	log.Info("[test] ready to login")
	param := make(map[string]string)
	param["loginEmail"] = "admin"
	param["password"] = "admin123"
	loginUrl := "http://azkaban_host/userReq/doLogin"
	// todo: verify login and get session id
	retBytes := http.DoPostFormRequest(loginUrl, param)
	log.Info("[test] login azkaban ret: %s", string(retBytes))
}
