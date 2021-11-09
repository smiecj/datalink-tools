package main

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

const (
	testAzkabanAddress = "azkaban_address"
	testUserName       = "admin"
	testPassword       = "admin"
)

func TestLogin(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testAzkabanAddress, Username: testUserName, Password: testPassword})
	err := client.login()
	require.Equal(t, err, nil)
}

func TestGetMedia(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testAzkabanAddress, Username: testUserName, Password: testPassword})
	err := client.login()
	mediaArr, _ := client.GetMedias()
	log.Info("[TestGetMedia] get media len: %d", len(mediaArr))
	require.Equal(t, err, nil)
}

func TestGetTasks(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testAzkabanAddress, Username: testUserName, Password: testPassword})
	err := client.login()
	require.Equal(t, err, nil)

	taskArr, err := client.GetTasks()
	log.Info("[TestGetTasks] get tasks len: %d", len(taskArr))
	require.Equal(t, err, nil)
}

func TestGetMappings(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testAzkabanAddress, Username: testUserName, Password: testPassword})
	err := client.login()
	mediaArr, _ := client.GetMappings()
	log.Info("[TestGetMappings] get mapping len: %d", len(mediaArr))
	require.Equal(t, err, nil)
}
