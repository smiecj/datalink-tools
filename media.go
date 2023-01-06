package tools

import (
	"fmt"
	"strings"
)

/*
// 查询 RDB Media 返回信息
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
type RDBStatusConfig struct {
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
type RDBConnectConfig struct {
	Encoding        string `json:"encoding"`
	MediaSourceType string `json:"mediaSourceType"`
	Name            string `json:"name"`
	Port            int    `json:"port"`
	ReadConfig      struct {
		Hosts    []string `json:"hosts"`
		Username string   `json:"username"`
		Password string   `json:"password"`
	} `json:"readConfig"`
	WriteConfig struct {
		Username  string `json:"username"`
		WriteHost string `json:"writeHost"`
	} `json:"writeConfig"`
}

// 关系型数据库介质信息
type RDBMedia struct {
	RDBStatusConfig  RDBStatusConfig  `json:"basicDataSourceConfig"`
	RDBConnectConfig RDBConnectConfig `json:"rdbMediaSrcParameter"`
	Schema           string           `json:"schema"`
}

func (media RDBMedia) GetReadHost() string {
	return media.RDBConnectConfig.ReadConfig.Hosts[0]
}

func (media RDBMedia) GetReadPort() int {
	return media.RDBConnectConfig.Port
}

func (media RDBMedia) GetReadUserName() string {
	return media.RDBConnectConfig.ReadConfig.Username
}

func (media RDBMedia) GetReadPassword() string {
	return media.RDBConnectConfig.ReadConfig.Password
}

func (media RDBMedia) String() string {
	return fmt.Sprintf("[RDBMedia] name: %s, address: %s:%s@%s:%d/%s", media.RDBConnectConfig.Name,
		media.RDBConnectConfig.ReadConfig.Username, media.RDBConnectConfig.ReadConfig.Password, media.RDBConnectConfig.ReadConfig.Hosts[0], media.RDBConnectConfig.Port,
		media.Schema)
}

// 查询关系型数据库介质列表接口返回信息 (/mediaSource/initMediaSource)
type QueryRDBMediaListRet struct {
	MediaList       []RDBMedia `json:"aaData"`
	Length          int        `json:"length"`
	PageNum         int        `json:"pageNum"`
	PageSize        int        `json:"pageSize"`
	Pages           int        `json:"pages"`
	RecordsFiltered int        `json:"recordsFiltered"`
	RecordsTotal    int        `json:"recordsTotal"`
	Size            int        `json:"size"`
	Start           int        `json:"start"`
}

/*
// 查询 kudu media 返回信息
{
  "aaData": [
    {
      "createTime": 1636647404000,
      "desc": "desc",
      "id": 4,
      "kuduMediaSrcParameter": {
        "bufferSize": 0,
        "database": "impala::db_name",
        "host2Ports": [
          "kudu_master_host:kudu_master_port"
        ],
        "impalaCconfigs": [
          {
            "host": "",
            "port": ""
          }
        ],
        "kuduMasterConfigs": [
          {
            "host": "kudu_master_host",
            "port": kudu_master_port
          }
        ],
        "mediaSourceType": "KUDU"
      },
      "name": "kudu_test_kudu"
    }
  ],
  "draw": 0,
  "length": 0,
  "pageNum": 0,
  "pageSize": 10,
  "pages": 0,
  "recordsFiltered": 0,
  "recordsTotal": 0,
  "size": 0,
  "start": 0
}
*/

// 查询kudu数据库介质列表接口返回信息 (/kudu/initKudu)
type QueryKuduMediaListRet struct {
	MediaList       []KuduMedia `json:"aaData"`
	Length          int         `json:"length"`
	PageNum         int         `json:"pageNum"`
	PageSize        int         `json:"pageSize"`
	Pages           int         `json:"pages"`
	RecordsFiltered int         `json:"recordsFiltered"`
	RecordsTotal    int         `json:"recordsTotal"`
	Size            int         `json:"size"`
	Start           int         `json:"start"`
}

// kudu media
type KuduMedia struct {
	Id                    int               `json:"int"`
	Name                  string            `json:"name"`
	CreateTime            int               `json:"createTime"`
	Desc                  string            `json:"desc"`
	KuduMediaSrcParameter KuduConnectConfig `json:"kuduMediaSrcParameter"`
}

func (media KuduMedia) Database() string {
	// 特殊逻辑: datalink 借助 impala 客户端同步数据，表名前面需要加上 impala::
	return strings.ReplaceAll(media.KuduMediaSrcParameter.Database, "impala::", "")
}

func (media KuduMedia) String() string {
	return fmt.Sprintf("[KuduMedia] name: %s, address: %s", media.Name, media.KuduMediaSrcParameter.Host2Ports)
}

// kudu 连接配置
type KuduConnectConfig struct {
	BufferSize      int      `json:"bufferSize"`
	Database        string   `json:"database"`
	Host2Ports      []string `json:"host2Ports"`
	MediaSourceType string   `json:"mediaSourceType"`

	ImpalaCconfigs []struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"impalaCconfigs"`
	KuduMasterConfigs []struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"kuduMasterConfigs"`
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
	defaultQueryRDBMediaStrFormat = `{"draw":1,"columns":[{"data":"id","name":"","searchable":true,"orderable":true,
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

	// 查询kudu参数：缺省，接口为GET 不需要传参
	defaultQueryKuduMediaStrFormat = ""

	// 默认保存RDB介质参数
	// 需要进行format 的参数:
	// {media_name}、{db_write_host}、{db_read_host}、{db_write_user}、{db_read_user}
	// {db_write_password}、{db_read_password}
	defaultSaveRDBMediaStrFormat = `rdbMediaSrcParameter.name={media_name}&rdbMediaSrcParameter.namespace=tmp&
  rdbMediaSrcParameter.mediaSourceType=MYSQL&rdbMediaSrcParameter.encoding=utf-8&rdbMediaSrcParameter.port=3306&
  rdbMediaSrcParameter.writeConfig.writeHost={db_write_host}&rdbMediaSrcParameter.writeConfig.username={db_write_user}&
  mysqlWritePsw={db_write_password}&mysqlReadPsw={db_read_password}&sqlserverWritePsw=&sqlserverReadPsw=&
  rdbMediaSrcParameter.writeConfig.password={db_write_password}&rdbMediaSrcParameter.readConfig.hosts%5B0%5D={db_read_host}&
  rdbMediaSrcParameter.readConfig.username={db_read_user}&rdbMediaSrcParameter.readConfig.password={db_read_password}&
  rdbMediaSrcParameter.readConfig.hosts%5B1%5D={db_read_host}&rdbMediaSrcParameter.readConfig.etlHost={db_read_host}&
  xxx.etlUserName=&xxx.etlPassWord=&rdbMediaSrcParameter.desc=test&basicDataSourceConfig.maxWait=60000&
  basicDataSourceConfig.minIdle=1&basicDataSourceConfig.initialSize=1&basicDataSourceConfig.maxActive=32&
  basicDataSourceConfig.maxIdle=32&basicDataSourceConfig.numTestsPerEvictionRun=-1&bas
  icDataSourceConfig.timeBetweenEvictionRunsMillis=60000&basicDataSourceConfig.removeAbandonedTimeout=300&
  basicDataSourceConfig.minEvictableIdleTimeMillis=3000000`

	// 默认保存kudu 介质参数
	// 需要进行format 的参数: {kudu_master_host}、{kudu_master_port}、{name}、{desc}、{db_name}
	defaultSaveKuduMediaStrFormat = `name={name}&kuduMediaSrcParameter.kuduMasterConfigs%5B0%5D.host={kudu_master_host}&
  kuduMediaSrcParameter.kuduMasterConfigs%5B0%5D.port={kudu_master_port}&kuduMediaSrcParameter.impalaCconfigs%5B0%5D.host=&
  kuduMediaSrcParameter.impalaCconfigs%5B0%5D.port=&kuduMediaSrcParameter.database={db_name}&desc={desc}`
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
type QueryRDBMediaParam struct {
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

// 生成 备份介质需要的基本信息
// 这里只会包含最关键的配置，一些默认配置不会体现
type RDBMediaBackupConfig struct {
	Name   string
	DBPort int

	// 读写节点配置
	// ETL-Host 和 WriteHost 保持一致
	DBWriteHost     string
	DBWriteUser     string
	DBWritePassword string
	DBReadHost      string
	DBReadUser      string
	DBReadPassword  string
}

// 从 rdb media 配置中生成 rdb media backup
// 数据库密码统一配置
func BuildRDBMediaBackupConf(media RDBMedia, password string) RDBMediaBackupConfig {
	conf := RDBMediaBackupConfig{
		Name:            media.RDBConnectConfig.Name,
		DBPort:          media.RDBConnectConfig.Port,
		DBWriteHost:     media.RDBConnectConfig.WriteConfig.WriteHost,
		DBWriteUser:     media.RDBConnectConfig.WriteConfig.Username,
		DBWritePassword: password,
		DBReadHost:      media.RDBConnectConfig.ReadConfig.Hosts[0],
		DBReadUser:      media.RDBConnectConfig.ReadConfig.Username,
		DBReadPassword:  password,
	}
	return conf
}

// kudu 备份配置
type KuduMediaBackupConfig struct {
	Name           string
	Desc           string
	KuduMasterHost string
	KuduMasterPort int
	DBName         string
}

// 从kudu 介质信息，生成 kudu 介质备份信息
func BuildKuduMediaBackupConf(media KuduMedia, password string) KuduMediaBackupConfig {
	conf := KuduMediaBackupConfig{
		Name:           media.Name,
		Desc:           media.Desc,
		KuduMasterHost: media.KuduMediaSrcParameter.KuduMasterConfigs[0].Host,
		KuduMasterPort: media.KuduMediaSrcParameter.KuduMasterConfigs[0].Port,
		DBName:         media.KuduMediaSrcParameter.Database,
	}
	return conf
}

// 查询介质详细信息参数
/* type QueryMediaDetailParam struct {
} */
// 不需要整理，toEdit 接口会返回一整个 html，不方便解析
// 存入介质的时候，手动输入密码
// 或者后续提供一个根据 mysql host name 查询数据库密码的接口
