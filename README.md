## Build


- `git clone https://github.com/abhirockzz/accs-go-kafka-mysql.git`
- `cd accs-go-kafka-mysql`
- `zip accs-go-kafka-mysql.zip gokafkamysql.go start.sh manifest.json`

## Run locally

- Make sure `Kafka` and `MySQL` are ready (locally or in the cloud)
- set the following environment variables - `MYSQLCS_USER_NAME`, `MYSQLCS_USER_PASSWORD`, `MYSQLCS_CONNECT_STRING` (format `<host>:<port>/<dbName>`), `OEHCS_EXTERNAL_CONNECT_STRING` (`<host>:<port>`) and `OEHCS_TOPIC` (Kafka topic)
- `go get github.com/bsm/sarama-cluster`
- `go get github.com/go-sql-driver/mysql`
- `go run gokafkamysql.go`
 
## Deploy to Oracle Application Container Cloud

Check out the blog - ['Go' for Kafka on Oracle Cloud]()
