package main

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
	srcMediaSourceName string `json:"srcMediaSourceName"`
	srcMediaNamespace  string `json:"srcMediaNamespace"`
	srcMediaName       string `json:"srcMediaName"`
	srcMediaSourceId   int    `json:"srcMediaSourceId"`
	srcMediaId         int    `json:"srcMediaId"`
	// 以下四个参数分别对应: 目标介质名、目标库名、目标表名、目标介质id
	targetMediaSourceName string `json:"targetMediaSourceName"`
	targetMediaNamespace  string `json:"targetMediaNamespace"`
	targetMediaName       string `json:"targetMediaName"`
	targetMediaSourceId   int    `json:"targetMediaSourceId"`
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

// 查询参数和storage共用
