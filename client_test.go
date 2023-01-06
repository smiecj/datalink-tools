package tools

import (
	"testing"
	"time"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

func TestGetMedia(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	mediaArr, err := client.GetRDBMedias()
	log.Info("[TestGetMedia] get media len: %d", len(mediaArr))
	for _, currentMedia := range mediaArr {
		log.Info("[TestGetMedia] current media: %s", currentMedia)
	}
	require.Equal(t, err, nil)
}

func TestGetKuduMedia(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	mediaArr, err := client.GetKuduMedias()
	log.Info("[TestKuduMedia] get kudu media len: %d", len(mediaArr))
	for _, currentMedia := range mediaArr {
		log.Info("[TestKuduMedia] current kudu media: %s", currentMedia)
	}
	require.Equal(t, err, nil)
}

func TestGetTasks(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	taskArr, err := client.GetTasks()
	log.Info("[TestGetTasks] get tasks len: %d", len(taskArr))
	require.Equal(t, err, nil)
	for _, currentTask := range taskArr {
		log.Info("[TestGetTasks] current task: %s", currentTask)
	}
}

func TestGetMappings(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	mappingArr, err := client.GetMappings()
	log.Info("[TestGetMappings] get mapping len: %d", len(mappingArr))
	require.Equal(t, err, nil)
	for _, currentMapping := range mappingArr {
		log.Info("[TestGetMappings] current mapping: %s", currentMapping)
	}
}

func TestRestartTask(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	taskArr, _ := client.GetTasks()
	log.Info("[TestGetMappings] get task len: %d", len(taskArr))
	require.NotEmpty(t, taskArr)
	// 重启第一个任务
	task := taskArr[0]
	err := client.StopTask(task.Id)
	require.Empty(t, err)
	time.Sleep(5 * time.Second)
	err = client.StartTask(task.Id)
	require.Empty(t, err)
}

func TestGetTaskMapping(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	taskArr, _ := client.GetTasks()
	kuduMediaArr, _ := client.GetKuduMedias()
	rdbMediaArr, _ := client.GetRDBMedias()
	mappingArr, _ := client.GetMappings()
	fullMapping := GetTaskMapping(client, taskArr, mappingArr, rdbMediaArr, kuduMediaArr)
	for _, currentMapping := range fullMapping {
		log.Info("[test] map: %s", currentMapping)
		log.Info("[test] mapping: %s", currentMapping.Mapping)
	}
}

func TestUpdateTask(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	taskArr, _ := client.GetTasks()
	// 获取第一个任务，更新时间戳
	require.NotEmpty(t, taskArr)
	firstTask := taskArr[0]
	taskDetail, err := client.GetTask(firstTask.Id)
	require.Nil(t, err)
	log.Info("[test] get task detail: %s", taskDetail)
	taskDetail.setCurrentTimestamp()
	client.UpdateTask(taskDetail)
}

func TestRefresh(t *testing.T) {
	client := GetDataLinkClient(testConfigManager)
	err := client.Refresh()
	require.Nil(t, err)
}
