package tools

import "github.com/smiecj/go_common/util/log"

// 工具方法: 获取源表（如 mysql ）和目标表（如 kudu）的映射关系
func GetTaskMapping(client Client, taskArr []*Task, mappingArr []*Mapping, mediaArr []*RDBMedia, kuduMediaArr []*KuduMedia) (retArr []*TaskMapping) {
	// 通过任务名建立映射
	taskMap := make(map[string]*Task)
	for _, currentTask := range taskArr {
		taskMap[currentTask.Name] = currentTask
	}

	rdbMediaMap := make(map[string]*RDBMedia)
	for _, currentMedia := range mediaArr {
		rdbMediaMap[currentMedia.RDBConnectConfig.Name] = currentMedia
	}

	kuduMediaMap := make(map[string]*KuduMedia)
	for _, currentMedia := range kuduMediaArr {
		kuduMediaMap[currentMedia.Name] = currentMedia
	}

	for _, currentMapping := range mappingArr {
		if currentTask, currentRdbMedia, currentKuduMedia, ok := checkMap(currentMapping, taskMap, rdbMediaMap, kuduMediaMap); ok {
			newTaskMapping := TaskMapping{
				Task:      *currentTask,
				Mapping:   *currentMapping,
				RDBMedia:  *currentRdbMedia,
				KuduMedia: *currentKuduMedia,
				TaskName:  currentMapping.TaskName,
			}
			retArr = append(retArr, &newTaskMapping)
		} else {
			log.Warn("[GetTaskAndMapping] mapping: %s, cannot find task", currentMapping)
		}
	}
	return
}

// 检查映射和数据源的映射关系
func checkMap(mapping *Mapping, taskMap map[string]*Task, rdbMediaMap map[string]*RDBMedia, rdbKuduMediaMap map[string]*KuduMedia) (*Task, *RDBMedia, *KuduMedia, bool) {
	task, checkTaskRet := taskMap[mapping.TaskName]
	rdbMedia, checkSourceRet := rdbMediaMap[mapping.SrcMediaSourceName]
	kuduMedia, checkTargetRet := rdbKuduMediaMap[mapping.TargetMediaSourceName]
	return task, rdbMedia, kuduMedia, checkTaskRet && checkSourceRet && checkTargetRet
}
