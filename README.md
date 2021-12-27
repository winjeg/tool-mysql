# tool-mysql


[![CircleCI](https://circleci.com/gh/winjeg/tool-mysql/tree/circleci-project-setup.svg?style=svg)](https://circleci.com/gh/winjeg/tool-mysql/?branch=circleci-project-setup)
[![Go Report Card](https://goreportcard.com/badge/github.com/winjeg/tool-mysql)](https://goreportcard.com/report/github.com/winjeg/tool-mysql)
![Golang CI](https://golangci.com/badges/github.com/winjeg/tool-mysql.svg)
[![codecov](https://codecov.io/gh/winjeg/tool-mysql/branch/master/graph/badge.svg)](https://codecov.io/gh/winjeg/tool-mysql)

---
nothing more than a tool here.

## db-filler

db-filler is a program to fill database table with randomly generated data,
with the program you can fill a table rapidly and easily.

![demo](https://user-images.githubusercontent.com/7270177/55327676-2e121580-54bd-11e9-8def-9a5c23f73f16.gif)

## feature list
1. auto generate insert sql using multi go routines, user can specify how many sql to generate
2. save the generated insert sql to a file
3. execute a sql before everything begins to the target database
4. execute fill table operation using multi go routines
5. save the execution error into a user defined file
6. only do insert sql generation, and don't do insertion

## how to use

### Get it
```
go get github.com/winjeg/tool-mysql
```

### build it
```bash
git clone https://github.com/winjeg/db-filler
cd db-filler
go build .
```

### command line usage
use `db-filer -h` to view the help content below

```markdown
db-filler usage:
db-filler -H localhost -p 3306 -a 123456 -d test -t person -n 1000
======================================================
-H string
address of a mysql server
-a string
password/auth of the mysql server
-d string
database name  of the mysql server to operate on
-dwn int
the number of routines to perform ddl (default 10)
-e string
the sql to be executed before everything starts
-ef string
the file where to put the error encountered when executing sql, if nothing specified, the error will be only in stdout
-gwn int
the number of routines to generate sql  (default 10)
-iwn int
the number of routines to perform insertion (default 10)
-n int
the number of sql to be generated and inserted (default 10)
-og
only generate the sql, don't do insertion
-p int
port of the mysql server (default 3306)
-rn int
the number of rows in one sql (default 1)
-sf string
the file to save generated sql, if none specified, no sql will be saved
-t string
the table name to be used to generate sql and to be inserted
-u string
username of the mysql server
```
