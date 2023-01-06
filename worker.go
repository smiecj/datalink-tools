package tools

import (
	"fmt"
)

const (
	// 默认查询参数
	defaultQueryWorkerStrFormat = `{"draw":4,"columns":[{"data":"id","name":"","searchable":true,
	"orderable":true,"search":{"value":"","regex":false}},{"data":"workerName","name":"",
	"searchable":true,"orderable":true,"search":{"value":"","regex":false}},
	{"data":"workerState","name":"","searchable":true,"orderable":true,
	"search":{"value":"","regex":false}},{"data":"workerAddress","name":"","searchable":true,
	"orderable":true,"search":{"value":"","regex":false}},{"data":"restPort","name":"",
	"searchable":true,"orderable":true,"search":{"value":"","regex":false}},
	{"data":"groupName","name":"","searchable":true,"orderable":true,
	"search":{"value":"","regex":false}},{"data":"startTime","name":"","searchable":true,"orderable":true,
	"search":{"value":"","regex":false}},{"data":"createTime","name":"","searchable":true,"orderable":false,
	"search":{"value":"","regex":false}},{"data":"id","name":"","searchable":true,"orderable":false,
	"search":{"value":"","regex":false}}],"order":[{"column":0,"dir":"asc"}],
	"start":%d,"length":%d,"search":{"value":"","regex":false},"groupId":"-1"}`
)

// 查询worker列表参数
type queryWorkerParam struct {
	Draw         int           `json:"draw"`
	QueryColumns []QueryColumn `json:"columns"`
	Start        int           `json:"start"`
	Length       int           `json:"length"`
	TaskId       string        `json:"taskId"`
	Order        []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
}

// 获取 worker 接口返回
type queryWorkerRet struct {
	WorkerList   []worker `json:"aaData"`
	Length       int      `json:"length"`
	PageNum      int      `json:"pageNum"`
	PageSize     int      `json:"pageSize"`
	RecordsTotal int      `json:"recordsTotal"`
}

// worker 信息
type worker struct {
	Port int    `json:"restPort"`
	Host string `json:"workerAddress"`
	Name string `json:"workerName"`
}

func (ret queryWorkerRet) getWorkerAddressArr() (retArr []string) {
	for _, currentWorker := range ret.WorkerList {
		retArr = append(retArr, fmt.Sprintf("http://%s:%d", currentWorker.Host, currentWorker.Port))
	}
	return
}
