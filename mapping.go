package tools

import "fmt"

/*
{
  "aaData": [
    {
      "createTime": 1590210384000,
      "esUsePrefix": true,
      "geoPositionConf": "",
      "id": 2,
      "parameter": "{}",
      "skipIds": "",
      "srcMediaId": 2,
      "srcMediaName": "table_name",
      "srcMediaNamespace": "db_name",
      "srcMediaSourceId": 1,
      "srcMediaSourceName": "mysql_media_source_name",
      "targetMediaName": "table_name",
      "targetMediaNamespace": "db_name",
      "targetMediaSourceId": 2,
      "targetMediaSourceName": "mysql_media_target_name",
      "taskName": "mysql",
      "valid": true,
      "writePriority": 5
    }
  ],
  "draw": 3,
  "length": 10,
  "pageNum": 1,
  "pageSize": 10,
  "pages": 0,
  "recordsFiltered": 59,
  "recordsTotal": 59,
  "size": 0,
  "start": 0
}
*/

// 映射
type Mapping struct {
	Id int `json:"id"`
	// 以下五个参数分别对应: 源介质名、源库名、源表名、源介质id、源自增id（介质自增id >= 介质id）
	SrcMediaSourceName string `json:"srcMediaSourceName"`
	SrcMediaNamespace  string `json:"srcMediaNamespace"`
	SrcMediaName       string `json:"srcMediaName"`
	SrcMediaSourceId   int    `json:"srcMediaSourceId"`
	SrcMediaId         int    `json:"srcMediaId"`
	// 以下四个参数分别对应: 目标介质名、目标库名、目标表名、目标介质id
	TargetMediaSourceName string `json:"targetMediaSourceName"`
	TargetMediaNamespace  string `json:"targetMediaNamespace"`
	TargetMediaName       string `json:"targetMediaName"`
	TargetMediaSourceId   int    `json:"targetMediaSourceId"`

	// 其他参数
	TaskName string `json:"taskName"`
}

func (mapping Mapping) SourceTable() string {
	return mapping.SrcMediaName
}

func (mapping Mapping) TargetTable() string {
	return mapping.TargetMediaName
}

func (mapping Mapping) String() string {
	return fmt.Sprintf("[mapping] task: %s, src: %s, target: %s", mapping.TaskName, mapping.SrcMediaName, mapping.TargetMediaName)
}

// 查询任务列表接口返回信息 (/mediaMapping/initMediaMapping)
type QueryMappingListRet struct {
	MappingList     []Mapping `json:"aaData"`
	Length          int       `json:"length"`
	PageNum         int       `json:"pageNum"`
	PageSize        int       `json:"pageSize"`
	Pages           int       `json:"pages"`
	RecordsFiltered int       `json:"recordsFiltered"`
	RecordsTotal    int       `json:"recordsTotal"`
	Size            int       `json:"size"`
	Start           int       `json:"start"`
}

/******************** 以下是接口查询参数 *********************/

/*
{
  "draw": 3,
  "columns": [
    {
      "data": "id",
      "name": "",
      "searchable": true,
      "orderable": true,
      "search": {
        "value": "",
        "regex": false
      }
    }
  ],
  "order": [
    {
      "column": 4,
      "dir": "asc"
    }
  ],
  "start": 0,
  "length": 10,
  "search": {
    "value": "",
    "regex": false
  },
  "mediaSourceId": "-1",
  "targetMediaSourceId": "-1",
  "srcMediaName": "",
  "targetMediaName": "",
  "taskId": "-1"
}
*/

const (
	// 默认查询参数
	defaultQueryMappingStrFormat = `{"draw":3,"columns":[{"data":"id","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"taskName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"srcMediaSourceName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"srcMediaName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"targetMediaSourceName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"targetMediaName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"writePriority","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"valid","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"createTime","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":null,"name":"","searchable":false,"orderable":false,
  "search":{"value":"","regex":false}}],"order":[{"column":4,"dir":"asc"}],"start":%d,"length":%d,
  "search":{"value":"","regex":false},"mediaSourceId":"-1","targetMediaSourceId":"-1","srcMediaName":"",
  "targetMediaName":"","taskId":"-1"}`
)

// 查询映射列表参数
type QueryMappingParam struct {
	Draw                int           `json:"draw"`
	QueryColumns        []QueryColumn `json:"columns"`
	Start               int           `json:"start"`
	Length              int           `json:"length"`
	MediaSourceId       string        `json:"mediaSourceId"`
	TargetMediaSourceId string        `json:"targetMediaSourceId"`
	SrcMediaName        string        `json:"srcMediaName"`
	TargetMediaName     string        `json:"targetMediaName"`
	TaskId              string        `json:"taskId"`
	Order               []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
}
