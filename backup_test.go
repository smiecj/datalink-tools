package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestBackuper() *Backuper {
	backuper := GetBackuper(DatalinkOption{
		Address:  testDatalinkAddress,
		Username: testUserName,
		Password: testPassword,
	})
	return backuper
}

func TestBackupMedia(t *testing.T) {
	backuper := getTestBackuper()
	mediaCount, err := backuper.BackupMedia()
	require.Less(t, 0, mediaCount)
	require.Equal(t, nil, err)
}

func TestBackupTask(t *testing.T) {
	backuper := getTestBackuper()
	taskCount, err := backuper.BackupTask()
	require.Less(t, 0, taskCount)
	require.Equal(t, nil, err)
}

func TestBackupMapping(t *testing.T) {
	backuper := getTestBackuper()
	mappingCount, err := backuper.BackupMapping()
	require.Less(t, 0, mappingCount)
	require.Equal(t, nil, err)
}
