package main

const (
	// datalink api
	urlLogin        = "/userReq/doLogin"
	urlGetRDBMedia  = "/mediaSource/initMediaSource"
	urlGetKuduMedia = "/kudu/initKudu"
	urlGetTasks     = "/mysqlTask/mysqlTaskDatas"
	urlGetMapping   = "/mediaMapping/initMediaMapping"

	urlAddRDBMedia  = ""
	urlAddKuduMedia = "/kudu/doAdd"
)

const (
	// backup and syncer store/read from this folder
	localStorePath = "/tmp/datalink_backup"

	// db connector
	spaceDBName  = "datalink"
	spaceMedia   = "media"
	spaceTask    = "task"
	spaceMapping = "mapping"

	backupKey = "backup"
)

const (
	// for test
	testDatalinkAddress = "datalink_address"
	testUserName        = "admin"
	testPassword        = "admin"
)
