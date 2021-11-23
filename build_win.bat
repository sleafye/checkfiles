@echo off
:set GOPATH=%CD%;%GOPATH%

echo go build -o bin/checkfiles_win.exe main.go
go build -o bin/checkfiles_win.exe main.go