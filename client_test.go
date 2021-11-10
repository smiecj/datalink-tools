package main

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

func TestGetMedia(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testDatalinkAddress, Username: testUserName, Password: testPassword})
	mediaArr, err := client.GetMedias()
	log.Info("[TestGetMedia] get media len: %d", len(mediaArr))
	require.Equal(t, err, nil)
}

func TestGetTasks(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testDatalinkAddress, Username: testUserName, Password: testPassword})
	taskArr, err := client.GetTasks()
	log.Info("[TestGetTasks] get tasks len: %d", len(taskArr))
	require.Equal(t, err, nil)
}

func TestGetMappings(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: testDatalinkAddress, Username: testUserName, Password: testPassword})
	mediaArr, err := client.GetMappings()
	log.Info("[TestGetMappings] get mapping len: %d", len(mediaArr))
	require.Equal(t, err, nil)
}
