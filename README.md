# Minimal MySQL Query Client

A lightweight command-line MySQL query tool that only supports SELECT operations for data security.

## 1. Build
### 1.1 Build 
```shell
go build -o dbquery ./main.go
```
### 1.2 Build for Linux
```shell
GOOS=linux GOARCH=amd64 go build -o dbquery ./main.go
```

## Usage
### 2.1 Connect using DSN
```shell
./dbquery -dsn "user:password@tcp(host:port)/dbname"
```

### 2.2 Connect using parameters
```shell
./dbquery -h <host> -P <port> -u <user> -p <password> -d <database>
```
