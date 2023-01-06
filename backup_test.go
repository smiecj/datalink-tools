package tools

import (
	"testing"

	"github.com/smiecj/go_common/config"
	"github.com/smiecj/go_common/util/file"
	"github.com/stretchr/testify/require"
)

var (
	testConfigManager config.Manager
)

func init() {
	configFilePath := file.FindFilePath(defaultConfigFile)
	testConfigManager, _ = config.GetYamlConfigManager(configFilePath)
}

func getTestBackuper() *Backuper {
	backuper := GetBackuper(testConfigManager)
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
