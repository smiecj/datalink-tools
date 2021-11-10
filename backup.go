package main

import (
	"github.com/smiecj/go_common/db"
	"github.com/smiecj/go_common/util/log"
)

// 备份
type Backuper struct {
	option  DatalinkOption
	client  Client
	storage db.RDBConnector
}

func GetBackuper(option DatalinkOption) *Backuper {
	backuper := new(Backuper)
	backuper.option = option
	backuper.client = GetDataLinkClient(option)
	backuper.storage = db.GetLocalFileConnector(localStorePath)
	return backuper
}

// 备份介质
func (backuper *Backuper) BackupMedia() (int, error) {
	// 从指定环境，获取当前所有介质
	mediaArr, err := backuper.client.GetMedias()
	if nil != err {
		log.Error("[BackupMedia] 获取介质失败，请检查: %s", err.Error())
		return 0, err
	}
	mediaObjArr := make([]interface{}, len(mediaArr))
	for index := 0; index < len(mediaArr); index++ {
		mediaObjArr[index] = mediaArr[index]
	}

	// 保存备份的数据到存储中（因为两个环境可能网络不互通，所以备份的时候建议保存到文件中）
	ret, err := backuper.storage.Insert(db.InsertSetSpace(spaceDBName, spaceMedia), db.InsertAddObjectArr(mediaObjArr))
	if nil != err {
		log.Error("[BackupMedia] 备份介质失败，请检查: %s", err.Error())
		return 0, err
	}

	log.Info("[BackupMedia] 备份介质成功，总备份介质数: %d", ret.AffectedRows)
	return len(mediaArr), nil
}

// 备份任务
func (backuper *Backuper) BackupTask() (int, error) {
	taskArr, err := backuper.client.GetTasks()
	if nil != err {
		log.Error("[BackupTask] 获取任务失败，请检查: %s", err.Error())
		return 0, err
	}
	taskObjArr := make([]interface{}, len(taskArr))
	for index := 0; index < len(taskArr); index++ {
		taskObjArr[index] = taskArr[index]
	}

	ret, err := backuper.storage.Insert(db.InsertSetSpace(spaceDBName, spaceTask), db.InsertAddObjectArr(taskObjArr))
	if nil != err {
		log.Error("[BackupTask] 备份任务失败，请检查: %s", err.Error())
		return 0, err
	}

	log.Info("[BackupTask] 备份任务成功，总备份任务数: %d", ret.AffectedRows)
	return len(taskArr), nil
}

// 备份映射
func (backuper *Backuper) BackupMapping() (int, error) {
	mappingArr, err := backuper.client.GetMappings()
	if nil != err {
		log.Error("[BackupMapping] 获取映射失败，请检查: %s", err.Error())
		return 0, err
	}
	mappingObjArr := make([]interface{}, len(mappingArr))
	for index := 0; index < len(mappingArr); index++ {
		mappingObjArr[index] = mappingArr[index]
	}

	ret, err := backuper.storage.Insert(db.InsertSetSpace(spaceDBName, spaceMapping), db.InsertAddObjectArr(mappingObjArr))
	if nil != err {
		log.Error("[BackupMapping] 备份映射失败，请检查: %s", err.Error())
		return 0, err
	}

	log.Info("[BackupMapping] 备份映射成功，总备份映射数: %d", ret.AffectedRows)
	return len(mappingArr), nil
}
