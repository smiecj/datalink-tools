package tools

import (
	"fmt"
)

const (
	defaultConfigFile = "config_local.yml"

	datalinkSpaceName = "datalink"

	// datalink api
	urlLogin        = "/userReq/doLogin"
	urlRDBMedia     = "/mediaSource/initMediaSource"
	urlKuduMedia    = "/kudu/initKudu"
	urlTasks        = "/mysqlTask/mysqlTaskDatas"
	urlGetTask      = "/mysqlTask/getMysqlTask"
	urlUpdateTask   = "/mysqlTask/doUpdateMysqlTask"
	urlMapping      = "/mediaMapping/initMediaMapping"
	urlAddRDBMedia  = ""
	urlAddKuduMedia = "/kudu/doAdd"
	urlStopTask     = "/mysqlTask/pauseMysqlTask"
	urlStartTask    = "/mysqlTask/resumeMysqlTask"
	urlGetWorker    = "/worker/initWorker"
	urlRefreshKudu  = "/flush/reloadKudu"
)

type datalinkConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	SessionInternal int    `yaml:"session_internal"`
}

func (conf *datalinkConfig) address() string {
	return fmt.Sprintf("http://%s:%s", conf.Host, conf.Port)
}

func (conf *datalinkConfig) urlLogin() string {
	return fmt.Sprintf("%s%s", conf.address(), urlLogin)
}

func (conf *datalinkConfig) urlRDBMedia() string {
	return fmt.Sprintf("%s%s", conf.address(), urlRDBMedia)
}

func (conf *datalinkConfig) urlKuduMedia() string {
	return fmt.Sprintf("%s%s", conf.address(), urlKuduMedia)
}

func (conf *datalinkConfig) urlTasks() string {
	return fmt.Sprintf("%s%s", conf.address(), urlTasks)
}

func (conf *datalinkConfig) urlGetTask(id int) string {
	return fmt.Sprintf("%s%s?id=%d", conf.address(), urlGetTask, id)
}

func (conf *datalinkConfig) urlUpdateTask() string {
	return fmt.Sprintf("%s%s", conf.address(), urlUpdateTask)
}

func (conf *datalinkConfig) urlMapping() string {
	return fmt.Sprintf("%s%s", conf.address(), urlMapping)
}

func (conf *datalinkConfig) urlStopTask(id int) string {
	return fmt.Sprintf("%s%s?id=%d", conf.address(), urlStopTask, id)
}

func (conf *datalinkConfig) urlStartTask(id int) string {
	return fmt.Sprintf("%s%s?id=%d", conf.address(), urlStartTask, id)
}

func (conf *datalinkConfig) urlGetWorker() string {
	return fmt.Sprintf("%s%s", conf.address(), urlGetWorker)
}

func (conf *datalinkConfig) urlRefreshKuduMedia(taskNodeAddress string, kuduMediaId int) string {
	return fmt.Sprintf("%s%s/%d", taskNodeAddress, urlRefreshKudu, kuduMediaId)
}
