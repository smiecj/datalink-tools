media:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestGetMedia"

kudu_media:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestGetKuduMedia"

task:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestGetTasks"

mapping:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestGetMappings"

task_map:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestGetTaskMapping"

task_update:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestUpdateTask"

restart:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestRestartTask"

refresh:
	go test -count=1 -v github.com/smiecj/datalink-tools -run="TestRefresh"