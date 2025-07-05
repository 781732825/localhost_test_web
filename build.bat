@echo off

REM 设置Golang环境变量
set BINARY_NAME=host
set BUILD_DIR=build

if not exist %BUILD_DIR% (
  mkdir %BUILD_DIR%
)

REM 编译Windows版本
set GOOS=windows
set GOARCH=386
echo 正在编译Windows x86版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_windows_x86.exe

set GOARCH=amd64
echo 正在编译Windows x64版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_windows_x64.exe

REM 编译Linux版本
set GOOS=linux
set GOARCH=386
echo 正在编译Linux x86版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_x86

set GOOS=linux
set GOARCH=amd64
echo 正在编译Linux x64版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_x64

set GOOS=linux
set GOARCH=mips
echo 正在编译Linux mips 版本(MIPS架构)...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_mips

set GOOS=linux
set GOARCH=mipsle
echo 正在编译Linux mips 版本(小端版MIPS架构)...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_mipsle

set GOOS=linux
set GOARCH=arm
echo 正在编译Linux arm x32版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_arm


set GOOS=linux
set GOARCH=arm64
echo 正在编译Linux arm x64版本...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_arm64


pause
