:@echo off
:set GOPATH=%CD%;%GOPATH%

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
echo go build -o bin/checkfiles_mac main.go
go build -o bin/checkfiles_mac main.go