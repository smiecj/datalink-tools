// package main: datalink tool: datalink 任务备份、同步工具
package main

import (
	"flag"

	"github.com/smiecj/go_common/util/log"
)

const (
	CommandShow   = "show"
	CommandBackup = "backup"
	CommandSync   = "sync"
)

var (
	command  = *flag.String("command", "show", "Command")
	address  = *flag.String("address", "http://localhost:18888", "datalink address")
	username = *flag.String("username", "admin", "datalink login username")
	password = *flag.String("password", "admin123", "datalink login password")
)

func main() {
	switch command {
	case CommandShow:
		log.Info("This is datalink tools!")
	case CommandBackup:
		backuper := GetBackuper(DatalinkOption{
			Address:  address,
			Username: username,
			Password: password,
		})
		mediaBackupCount, _ := backuper.BackupMedia()
		taskBackupCount, _ := backuper.BackupTask()
		mappingBackupCount, _ := backuper.BackupMapping()

		log.Info("[main] backup datalink config result: media backup count: %d, task backup count: %d, "+
			"mapping backup count: %d", mediaBackupCount, taskBackupCount, mappingBackupCount)
	case CommandSync:
		syncer := GetSyncer(DatalinkOption{
			Address:  address,
			Username: username,
			Password: password,
		})
		storageCount, err := syncer.SyncMedia()
		if nil != err {
			log.Error("[main] backup datalink storage config failed: %s", err.Error())
			return
		} else {
			log.Info("[main] backup datalink storage config success: %d", storageCount)
		}
	}

}
