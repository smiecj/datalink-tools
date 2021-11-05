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
		mediaCount, err := backuper.BackupMedia()
		if nil != err {
			log.Error("[main] backup datalink media config failed: %s", err.Error())
			return
		} else {
			log.Info("[main] backup datalink media config success: %d", mediaCount)
		}
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
