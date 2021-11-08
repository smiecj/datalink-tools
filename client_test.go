package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	client := GetDataLinkClient(DatalinkOption{Address: "azkaban_address", Username: "admin", Password: "admin"})
	err := client.login()
	require.Equal(t, err, nil)
}
