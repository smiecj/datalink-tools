package main

/*
{
  "aaData": [
    {
      "basicDataSourceConfig": {
        "initialSize": 1,
        "maxActive": 32,
        "maxIdle": 32,
        "maxWait": 60000,
        "minEvictableIdleTimeMillis": 300000,
        "minIdle": 1,
        "numTestsPerEvictionRun": -1,
        "removeAbandonedTimeout": 300,
        "timeBetweenEvictionRunsMillis": 60000
      },
      "createTime": 1590111677000,
      "id": 1,
      "rdbMediaSrcParameter": {
        "encoding": "UTF-8",
        "mediaSourceType": "MYSQL",
        "name": "my_test_mysql",
        "port": 0,
        "readConfig": {
          "hosts": [
            "mysql1.com",
            "mysql1.com"
          ],
          "username": "canal"
        },
        "writeConfig": {
          "username": "canal",
          "writeHost": "mysql1.host"
        }
      }
    }
  ],
  "draw": 1,
  "length": 10,
  "pageNum": 1,
  "pageSize": 10,
  "pages": 0,
  "recordsFiltered": 20,
  "recordsTotal": 20,
  "size": 0,
  "start": 0
}
*/

// 数据库状态配置，如连接超时时间
// 主要是 Druid 连接池相关配置，参考: https://www.jianshu.com/p/be9dbe640daf
type DBStatusConfig struct {
	// 初始化创建连接数
	InitialSize int `json:"initialSize"`
	MaxActive   int `json:"maxActive"`
	MaxIdle     int `json:"maxIdle"`
	MinIdle     int `json:"minIdle"`
	MaxWait     int `json:"maxWait"`
	// 超过 minIdle 的那部分连接的空闲超时时间
	MinEvictableIdleTimeMillis    int `json:"minEvictableIdleTimeMillis"`
	NumTestsPerEvictionRun        int `json:"numTestsPerEvictionRun"`
	RemoveAbandonedTimeout        int `json:"removeAbandonedTimeout"`
	TimeBetweenEvictionRunsMillis int `json:"timeBetweenEvictionRunsMillis"`
}

// 数据库连接配置，如数据源地址、连接账户
type DBConnectConfig struct {
	Encoding        string `json:"encoding"`
	MediaSourceType string `json:"mediaSourceType"`
	Name            string `json:"name"`
	Port            int    `json:"port"`
	ReadConfig      struct {
		Hosts    []string `json:"hosts"`
		Username string   `json:"username"`
	} `json:"readConfig"`
	WriteConfig struct {
		Username  string `json:"username"`
		WriteHost string `json:"writeHost"`
	} `json:"writeConfig"`
}

// 介质信息
type Media struct {
	DBStatusConfig  DBStatusConfig  `json:"basicDataSourceConfig"`
	DBConnectConfig DBConnectConfig `json:"rdbMediaSrcParameter"`
}

// 查询介质列表接口返回信息 (/mediaSource/initMediaSource)
type QueryMediaListRet struct {
	MediaList       []Media `json:"aaData"`
	Length          int     `json:"length"`
	PageNum         int     `json:"pageNum"`
	PageSize        int     `json:"pageSize"`
	Pages           int     `json:"pages"`
	RecordsFiltered int     `json:"recordsFiltered"`
	RecordsTotal    int     `json:"recordsTotal"`
	Size            int     `json:"size"`
	Start           int     `json:"start"`
}

/******************** 以下是接口查询参数 *********************/
/*
{
  "draw": 1,
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
    },
    {
      "data": "rdbMediaSrcParameter.name",
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
  "mediaSourceType": "-1",
  "name": "",
  "ip": ""
}
*/

const (
	// 默认查询参数
	defaultQueryMediaStrFormat = `{"draw":1,"columns":[{"data":"id","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"rdbMediaSrcParameter.name","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"rdbMediaSrcParameter.mediaSourceType","name":"","searchable":true,
  "orderable":true,"search":{"value":"","regex":false}},{"data":"rdbMediaSrcParameter.writeConfig.writeHost","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},
  {"data":"rdbMediaSrcParameter.writeConfig.username","name":"","searchable":true,"orderable":true,
  "search":{"value":"","regex":false}},{"data":"rdbMediaSrcParameter.readConfig.hosts","name":"","searchable":true,
  "orderable":true,"search":{"value":"","regex":false}},{"data":"rdbMediaSrcParameter.readConfig.username","name":"",
  "searchable":true,"orderable":true,"search":{"value":"","regex":false}},{"data":"createTime","name":"","searchable":true,
  "orderable":false,"search":{"value":"","regex":false}},{"data":"id","name":"","searchable":true,"orderable":false,
  "search":{"value":"","regex":false}}],"order":[{"column":0,"dir":"asc"}],"start":%d,"length":%d,
  "search":{"value":"","regex":false},"mediaSourceType":"-1","name":"","ip":""}`
)

// 查询介质，需要查询的列参数
type QueryColumn struct {
	Data       string `json:"data"`
	Name       string `json:"name"`
	Searchable bool   `json:"searchable"`
	Orderable  bool   `json:"orderable"`
	Search     struct {
		Value string `json:"value"`
		Regex bool   `json:"regex"`
	} `json:"search"`
}

// 查询介质列表参数
type QueryMediaParam struct {
	Draw            int           `json:"draw"`
	QueryColumns    []QueryColumn `json:"columns"`
	Start           int           `json:"start"`
	Length          int           `json:"length"`
	MediaSourceType string        `json:"mediaSourceType"`
	Order           []struct {
		Column int    `json:"column"`
		Dir    string `json:"dir"`
	} `json:"order"`
}

// 生成查询介质参数，需要设置的参数通过传参直接传入

// 查询介质详细信息参数
/* type QueryMediaDetailParam struct {
} */
// 不需要整理，toEdit 接口会返回一整个 html，不方便解析
// 存入介质的时候，手动输入密码
// 或者后续提供一个根据 mysql host name 查询数据库密码的接口
