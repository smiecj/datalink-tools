package main

import (
	"testing"

	http "github.com/smiecj/go_common/http"
	"github.com/smiecj/go_common/util/log"
)

func TestLoginByHttpClient(t *testing.T) {
	log.Info("[test] ready to login")
	loginUrl := "http://azkaban_address/userReq/doLogin"
	httpClient := http.GetHTTPClient()
	rsp, _ := httpClient.Do(http.Url(loginUrl),
		http.PostWithUrlEncode(),
		http.AddParam("loginEmail", "admin"),
		http.AddParam("password", "admin"))
	log.Info("[test] login azkaban ret: body: %s", string(rsp.Body))
	for key, value := range rsp.Header {
		log.Info("[test] login azkaban ret: header: %s -> %s", key, value)
	}
}
