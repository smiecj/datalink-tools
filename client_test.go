package main

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: "azkaban_address", Username: "admin", Password: "admin"})
	err := client.login()
	require.Equal(t, err, nil)
}

func TestGetMedia(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: "azkaban_address", Username: "admin", Password: "admin"})
	err := client.login()
	mediaArr, _ := client.GetMedias()
	log.Info("[TestGetMedia] get media len: %d", len(mediaArr))
	require.Equal(t, err, nil)
}
