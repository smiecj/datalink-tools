package main

// 同步datalink 配置
type Syncer struct {
	Option DatalinkOption
}

func GetSyncer(option DatalinkOption) *Syncer {
	syncer := new(Syncer)
	syncer.Option = option
	return syncer
}

// todo: 同步介质
func (syncer *Syncer) SyncMedia() (int, error) {

	return 0, nil
}

// todo: 同步任务
func (syncer *Syncer) SyncTask() (int, error) {
	return 0, nil
}

// todo: 同步映射
func (syncer *Syncer) SyncMapping() (int, error) {
	return 0, nil
}
