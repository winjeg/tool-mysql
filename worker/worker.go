package worker

import (
	"github.com/winjeg/go-commons/log"
	"github.com/winjeg/tool-mysql/schema"
	"github.com/winjeg/tool-mysql/store"

	"fmt"
	"sync"
)

var (
	logger = log.GetLogger(nil)
	db     = store.GetDb()
)

func logError(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// GenerateSql  generate sqls using threads.
// totalNum : total rows of data to be generated
// threadNum: how many thread used to generate the sql
// bulkNum: how many rows in one insert
// tableName: the table name to be inserted
// dbName: the database name to be inserted
func GenerateSql(totalNum, threadNum, bulkNum int, tableName, dbName string) []string {
	// to check if table exits, if not, panic and exit the program
	_, _, err := store.GetFromResult(db.Exec(fmt.Sprintf("show create table %s", tableName)))
	if err != nil {
		logger.Panic("the table you specified may not exits", err.Error())
	}
	tableDef := schema.TableService.GetTableDefinition(db, tableName, dbName)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(threadNum)

	genNum := totalNum
	perNum := genNum / threadNum / bulkNum
	restNum := genNum % (threadNum * bulkNum)

	result := make([]string, 0, genNum)
	var lock sync.Mutex

	for i := 0; i < threadNum; i++ {
		go func(n int) {
			defer waitGroup.Done()
			ts := make([]string, 0, perNum)
			if n == 0 {
				// first thread takes the reset and original
				for j := 0; j < restNum+perNum; j++ {
					sql := tableDef.GenerateInsert(bulkNum)
					ts = append(ts, sql)
				}
			} else {
				for j := 0; j < perNum; j++ {
					sql := tableDef.GenerateInsert(bulkNum)
					ts = append(ts, sql)
				}
			}
			// put it to the result slice safely
			if len(ts) > 0 {
				lock.Lock()
				result = append(result, ts...)
				lock.Unlock()
			}
		}(i)
	}
	waitGroup.Wait()
	return result
}

// ParallelExecute parallel execute sql in n routines
// sqls: generated sqls
// workerNum: num of workers to insert sql
// returns the sqls that got wrong
func ParallelExecute(sqls []string, workerNum int) []string {
	var waitGroup sync.WaitGroup
	sqlLen := len(sqls)
	perNum := sqlLen / workerNum
	waitGroup.Add(workerNum)
	errorSqls := make([]string, 0, sqlLen)
	var lock sync.Mutex
	for i := 0; i < workerNum; i++ {
		go func(n int) {
			defer waitGroup.Done()
			terrorSqls := make([]string, 0, perNum)
			tsqls := sqls[n*perNum : n*perNum+perNum]
			if n == workerNum-1 {
				tsqls = sqls[n*perNum:]
			}
			for _, v := range tsqls {
				_, _, err := store.GetFromResult(db.Exec(v))
				if err != nil {
					logError(err)
					terrorSqls = append(terrorSqls, v)
				}
			}
			if len(terrorSqls) > 0 {
				lock.Lock()
				errorSqls = append(errorSqls, terrorSqls...)
				lock.Unlock()
			}
		}(i)
	}
	waitGroup.Wait()
	return errorSqls
}
