@echo off

REM ����Golang��������
set BINARY_NAME=host
set BUILD_DIR=build

if not exist %BUILD_DIR% (
  mkdir %BUILD_DIR%
)

REM ����Windows�汾
set GOOS=windows
set GOARCH=386
echo ���ڱ���Windows x86�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_windows_x86.exe

set GOARCH=amd64
echo ���ڱ���Windows x64�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_windows_x64.exe

REM ����Linux�汾
set GOOS=linux
set GOARCH=386
echo ���ڱ���Linux x86�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_x86

set GOOS=linux
set GOARCH=amd64
echo ���ڱ���Linux x64�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_x64

set GOOS=linux
set GOARCH=mips
echo ���ڱ���Linux mips �汾(MIPS�ܹ�)...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_mips

set GOOS=linux
set GOARCH=mipsle
echo ���ڱ���Linux mips �汾(С�˰�MIPS�ܹ�)...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_mipsle

set GOOS=linux
set GOARCH=arm
echo ���ڱ���Linux arm x32�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_arm


set GOOS=linux
set GOARCH=arm64
echo ���ڱ���Linux arm x64�汾...
go build -v -o %BUILD_DIR%/%BINARY_NAME%_linux_arm64


pause
