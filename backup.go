package main

// 备份
type Backuper struct {
	Option DatalinkOption
}

func GetBackuper(option DatalinkOption) *Backuper {
	backuper := new(Backuper)
	backuper.Option = option
	return backuper
}

// todo: 备份介质
func (backuper *Backuper) BackupMedia() (int, error) {
	// 从指定环境，获取当前所有介质

	// 保存配置到本地存储中
	return 0, nil
}

// todo: 备份映射

// todo: 备份任务
