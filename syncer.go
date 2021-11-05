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

// todo: 备份介质
func (syncer *Syncer) SyncMedia() (int, error) {
	return 0, nil
}
