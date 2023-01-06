package tools

import (
	"fmt"
	"time"
)

const (
	kuduWriterParameterClass = "com.ucar.datalink.domain.plugin.writer.kudu.KuduWriterParameter"
)

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

// 任务 (列表参数)
type Task struct {
	Id                 int    `json:"id"`
	GroupId            int    `json:"groupId"`
	WorkerId           int    `json:"workerId"`
	Name               string `json:"taskName"`
	Description        string `json:"taskDesc"`
	SyncStatus         string `json:"taskSyncStatus"`
	ListenedState      string `json:"listenedState"`
	TargetState        string `json:"targetState"`
	CurrentLogFile     string `json:"currentLogFile"`
	CurrentLogPosition int    `json:"currentLogPosition"`
}

func (task Task) String() string {
	return fmt.Sprintf("[task] name: %s, status: %s", task.Name, task.ListenedState)
}

// 任务 (详细参数)
type TaskDetail struct {
	TaskBasicInfo TaskBasicInfo `json:"taskBasicInfo"`

	MysqlReaderParameter MysqlReaderParameter `json:"mysqlReaderParameter"`

	WriterParameterMap WriterParameterMap `json:"writerParameterMap"`
}

func (detail TaskDetail) String() string {
	return fmt.Sprintf("[TaskDetail] name: %s", detail.TaskBasicInfo.TaskName)
}

// 设置同步时间戳 为当前时间 (目前仅供测试)
func (detail *TaskDetail) setCurrentTimestamp() {
	detail.MysqlReaderParameter.StartTimeStamps = int(time.Now().UnixMilli())
}

// 设置同步时间戳 为指定时间
func (detail *TaskDetail) SetTimestamp(timestamp int) {
	detail.MysqlReaderParameter.StartTimeStamps = timestamp
}

// 任务详情和更新参数转换
func (detail *TaskDetail) getDetailUpdate() *TaskDetailUpdate {
	ret := new(TaskDetailUpdate)
	ret.MysqlReaderParameter = detail.MysqlReaderParameter
	ret.TaskBasicInfo = detail.TaskBasicInfo
	ret.WriterParameterMap = detail.WriterParameterMap
	ret.WriterParameterMap.WriterKudu.Type = kuduWriterParameterClass
	ret.WriterParameterMap.WriterKudu.PluginClass = ""
	// ret.MysqlReaderParameter.FilteredEventTypes = nil
	return ret
}

type TaskBasicInfo struct {
	Id          int    `json:"id"`
	TaskName    string `json:"taskName"`
	TaskDesc    string `json:"taskDesc"`
	TargetState string `json:"targetState"`
	GroupId     int    `json:"groupId"`
}

type MysqlReaderParameter struct {
	MediaSourceId                      int      `json:"mediaSourceId"`
	StartTimeStamps                    int      `json:"startTimeStamps"`
	BlackFilter                        string   `json:"blackFilter"`
	BatchTimeout                       int      `json:"batchTimeout"`
	DefaultConnectionTimeoutInSeconds  int      `json:"defaultConnectionTimeoutInSeconds"`
	DetectingIntervalInSeconds         int      `json:"detectingIntervalInSeconds"`
	DetectingRetryTimes                int      `json:"detectingRetryTimes"`
	DetectingSQL                       string   `json:"detectingSQL"`
	DetectingTimeoutThresholdInSeconds int      `json:"detectingTimeoutThresholdInSeconds"`
	Dump                               bool     `json:"dump"`
	DumpDetail                         bool     `json:"dumpDetail"`
	FallbackIntervalInSeconds          int      `json:"fallbackIntervalInSeconds"`
	FilteredEventTypes                 []string `json:"filteredEventTypes"`
	MemoryStorageBufferMemUnit         int      `json:"memoryStorageBufferMemUnit"`
	MemoryStorageBufferSize            int      `json:"memoryStorageBufferSize"`
	MessageBatchSize                   int      `json:"messageBatchSize"`
	ReceiveBufferSize                  int      `json:"receiveBufferSize"`
	SendBufferSize                     int      `json:"sendBufferSize"`
	GroupSinkMode                      string   `json:"groupSinkMode"`
	DdlSync                            bool     `json:"ddlSync"`
	PerfStatistic                      bool     `json:"perfStatistic"`
	Parallel                           bool     `json:"parallel"`
	GtidEnable                         bool     `json:"gtidEnable"`
	MultiplexingRead                   bool     `json:"multiplexingRead"`
}

type WriterParameterMap struct {
	WriterKudu struct {
		PluginClass string `json:"pluginClass"`
		Type        string `json:"@type"`
		PoolSize    int    `json:"poolSize"`
		BatchSize   int    `json:"batchSize"`
	} `json:"writer-kudu"`
}

// 任务详细参数（更新任务接口用）
type TaskDetailUpdate struct {
	TaskBasicInfo TaskBasicInfo `json:"taskBasicInfo"`

	MysqlReaderParameter MysqlReaderParameter `json:"mysqlReaderParameter"`

	WriterParameterMap WriterParameterMap `json:"writerParameterMap"`
}

// 查询任务列表接口返回信息 (/mysqlTask/mysqlTaskDatas)
type QueryTaskListRet struct {
	TaskList        []Task `json:"aaData"`
	Length          int    `json:"length"`
	PageNum         int    `json:"pageNum"`
	PageSize        int    `json:"pageSize"`
	Pages           int    `json:"pages"`
	RecordsFiltered int    `json:"recordsFiltered"`
	RecordsTotal    int    `json:"recordsTotal"`
	Size            int    `json:"size"`
	Start           int    `json:"start"`
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

const (
	// 默认查询参数
	defaultQueryTaskStrFormat = `{"draw":2,"columns":[{"data":"id","name":"","searchable":true,"orderable":false,
  "search":{"value":"","regex":false}},{"data":"id","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"taskName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"detail","name":"","searchable":true,"orderable":false,
  "search":{"value":"","regex":false}},{"data":"targetState","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"listenedState","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"groupId","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"workerId","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"currentTimeStamp","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"startTime","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"readerIp","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"latestEffectSyncLogFileName","name":"","searchable":true,
  "orderable":true,"search":{"value":"","regex":false}},{"data":"latestEffectSyncLogFileOffset","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"taskSyncStatus","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"shadowCurrentTimeStamp","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},
  {"data":"shadowLatestEffectSyncLogFileName","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"shadowLatestEffectSyncLogFileOffset","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":null,"name":"",
  "searchable":false,"orderable":false,"search":{"value":"","regex":false}}],"order":[{"column":0,"dir":"asc"}],
  "start":%d,"length":%d,"search":{"value":"","regex":false},"readerMediaSourceId":"-1","groupId":"-1","id":"-1"}`
)

// 查询任务列表参数
type QueryTaskParam struct {
	Draw                int           `json:"draw"`
	QueryColumns        []QueryColumn `json:"columns"`
	Start               int           `json:"start"`
	Length              int           `json:"length"`
	ReaderMediaSourceId string        `json:"readerMediaSourceId"`
	GroupId             string        `json:"groupId"`
	Id                  string        `json:"id"`
	Order               []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
}
