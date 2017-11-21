#!/bin/sh
#enable proxy
#export HTTP_PROXY=http://www-proxy.us.oracle.com:80
#export HTTPS_PROXY=https://www-proxy.us.oracle.com:80

go get github.com/bsm/sarama-cluster
go get github.com/go-sql-driver/mysql
go run gokafkamysql.go