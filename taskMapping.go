package tools

import "fmt"

// 任务和映射的对应关系
// 默认存入 mysql 源表和 kudu 表的映射
type TaskMapping struct {
	Task      Task
	Mapping   Mapping
	RDBMedia  RDBMedia
	KuduMedia KuduMedia
	TaskName  string
}

func (mapping TaskMapping) FullSourceTable() string {
	return fmt.Sprintf("%s.%s", mapping.RDBMedia.Schema, mapping.Mapping.SourceTable())
}

func (mapping TaskMapping) FullTargetTable() string {
	return fmt.Sprintf("%s.%s", mapping.KuduMedia.Database(), mapping.Mapping.TargetTable())
}

func (mapping TaskMapping) String() string {
	return fmt.Sprintf("src media: %s, src table: %s, target: %s, target table: %s",
		mapping.RDBMedia, mapping.FullSourceTable(),
		mapping.KuduMedia, mapping.FullTargetTable())
}
