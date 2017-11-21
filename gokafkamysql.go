package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	mysqlUser := os.Getenv("MYSQLCS_USER_NAME")
	mysqlPwd := os.Getenv("MYSQLCS_USER_PASSWORD")
	mysqlConnString := os.Getenv("MYSQLCS_CONNECT_STRING")
	
	mysqlConnStringAndDB := strings.Split(mysqlConnString, "/")
    hostport, dbname := mysqlConnStringAndDB[0], mysqlConnStringAndDB[1]
    fmt.Println(hostport, dbname)

	mysqlDSNForDriver := mysqlUser + ":" + mysqlPwd + "@tcp(" + hostport + ")/" + dbname

	//db, err := sql.Open("mysql", "root:root@tcp(192.168.99.100:3306)/mysql")
	db, err := sql.Open("mysql", mysqlDSNForDriver)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	
	fmt.Println("connected to MySQL...")

	stmtIns, err := db.Prepare("INSERT INTO datadump (`topic`,`partition`,`offset`,`key`,`value`,`processedby`,`createdat`) VALUES(?,?,?,?,?,?,?)") 

	if err != nil {
		panic(err.Error()) 
	}
	defer stmtIns.Close() 

	ehcsBroker := os.Getenv("OEHCS_EXTERNAL_CONNECT_STRING")
	if ehcsBroker == "" {
		ehcsBroker = "192.168.99.100:9092"
	}

	ehcsTopic := os.Getenv("OEHCS_TOPIC")
	if ehcsTopic == "" {
		ehcsTopic = "test"
	}

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	brokers := []string{ehcsBroker}
	topics := []string{ehcsTopic}
	fmt.Println("Connecting to EHCS cluster at " + ehcsBroker)

	consumer, err := cluster.NewConsumer(brokers, "test-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to Kafka...")
	defer consumer.Close()
	
	//get ACCS app instance name
	appName := os.Getenv("ORA_APP_NAME")
	if appName == "" {
		appName = "accsgokafkamysql"
	}
	
	appInstance := os.Getenv("ORA_INSTANCE_NAME")
	if appInstance == "" {
		appInstance = "instance1"
	}
	
	processedby := appName + "_" + appInstance
	
	fmt.Println("Processed by instance "+ processedby)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

				_, err = stmtIns.Exec(msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value, processedby, time.Now().Local())
				if err != nil {
					panic(err.Error())
				} else {
					fmt.Println("Detailed pushed to MySQL")
				}

				consumer.MarkOffset(msg, "")
			}
		case <-signals:
			return
		}
	}

}
