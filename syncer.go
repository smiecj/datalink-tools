package main

// 同步datalink 配置
type Syncer struct {
	Option DatalinkOption
}

func GetSyncer(option DatalinkOption) *Backuper {
	backuper := new(Backuper)
	backuper.Option = option
	return backuper
}

// todo: 备份介质
func (backuper *Backuper) BackupMedia() (int, error) {

}
