@echo off

echo Start running...

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GOPROXY=https://goproxy.cn
go build main.go

echo End of running.

pause