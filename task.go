package main

/*
{
  "aaData": [
    {
      "currentLogFile": "mysql-bin.000733",
      "currentLogPosition": 188013952,
      "currentTimeStamp": "2021-11-03 17:37:03 \u00281635932223000\u0029",
      "detail": "",
      "groupId": 3,
      "id": 1,
      "latestEffectSyncLogFileName": "",
      "latestEffectSyncLogFileOffset": "",
      "listenedState": "RUNNING",
      "readerIp": "10.10.100.242",
      "shadowCurrentTimeStamp": "",
      "shadowLatestEffectSyncLogFileName": "",
      "shadowLatestEffectSyncLogFileOffset": "",
      "startTime": "2021-11-03 16:02:28",
      "targetState": "STARTED",
      "taskDesc": "hive_to_mysql",
      "taskName": "hive_to_mysql",
      "taskSyncStatus": "Init",
      "workerId": 3
    }
  ],
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

// 任务
type Task struct {
	id                 int    `json:"id"`
	groupId            int    `json:"groupId"`
	workerId           string `json:"workerId"`
	Name               string `json:"taskName"`
	Description        string `json:"taskDesc"`
	SyncStatus         string `json:"taskSyncStatus"`
	targetState        string `json:"targetState"`
	currentLogFile     string `json:"currentLogFile"`
	currentLogPosition string `json:"taskDesc"`
}

// 查询任务列表接口返回信息 (/mysqlTask/mysqlTaskDatas)
type QueryTaskListRet struct {
	TaskList        []Media `json:"aaData"`
	Length          int     `json:"length"`
	PageNum         int     `json:"pageNum"`
	PageSize        int     `json:"pageSize"`
	Pages           int     `json:"pages"`
	recordsFiltered int     `json:"recordsFiltered"`
	recordsTotal    int     `json:"recordsTotal"`
	Size            int     `json:"size"`
	Start           int     `json:"start"`
}

/******************** 以下是接口查询参数 *********************/
/*
{
  "draw": 2,
  "columns": [
    {
      "data": "id",
      "name": "",
      "searchable": true,
      "orderable": false,
      "search": {
        "value": "",
        "regex": false
      }
    }
  ],
  "order": [
    {
      "column": 0,
      "dir": "asc"
    }
  ],
  "start": 0,
  "length": 10,
  "search": {
    "value": "",
    "regex": false
  },
  "readerMediaSourceId": "-1",
  "groupId": "-1",
  "id": "-1"
}
*/

// 查询参数和storage共用
