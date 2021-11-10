package main

const (
	// datalink api
	urlLogin      = "/userReq/doLogin"
	urlGetMedia   = "/mediaSource/initMediaSource"
	urlGetTasks   = "/mysqlTask/mysqlTaskDatas"
	urlGetMapping = "/mediaMapping/initMediaMapping"
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
