package schema

import (
	"github.com/stretchr/testify/assert"
	"github.com/winjeg/tool-mysql/store"

	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var (
	db  = store.GetDb()
	tsc = TableService.GetTableDefinition(db, "label", "test")
)

func BenchmarkTableDefinition_GenerateInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tsc.GenerateInsert(1)
	}
	b.ReportAllocs()
}

func TestTbIdxes_Len(t *testing.T) {
	start := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		tsc.GenerateInsert(1)
	}
	fmt.Println(time.Now().UnixNano() - start)
}

const testCreateSql = `
	 CREATE TABLE %s (
		  id int(11) NOT NULL COMMENT 'the primary key',
		  name varchar(64) NOT NULL DEFAULT 'unknown' COMMENT 'name',
		  age tinyint(4) NOT NULL DEFAULT '0' COMMENT 'age',
		  deleted bit(1) NOT NULL DEFAULT b'0',
		  note varchar(255) DEFAULT NULL COMMENT 'desc',
		  created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
		  updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last updated',
		  PRIMARY KEY (id),
		  UNIQUE KEY idx_id (id),
		  KEY idx_name (name),
		  KEY idx_name_deleted (name) USING BTREE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='table for label information'
	`

func TestSchemaService(t *testing.T) {
	_, err := db.Exec(fmt.Sprintf(testCreateSql, "user_info"))
	assert.Nil(t, err)
	td := TableService.GetTableDefinition(db, "user_info", "test")
	td.GenerateInsert(10)
	_, err = json.Marshal(td)
	assert.Nil(t, err)
	originalName := td.TableName
	td.TableName += "_test"
	td2 := TableDefinition(td)
	td2.Indexes = nil
	td.GetCreateTableSql()
	sql := td.GetCreateTableSql()
	_, _, err = store.GetFromResult(db.Exec(sql))
	assert.Nil(t, err)
	_, _, err = store.GetFromResult(db.Exec(fmt.Sprintf("DROP TABLE %s", td.TableName)))
	assert.Nil(t, err)
	td.TableName = originalName

	// dropping cols and indexes
	lastIndex := td.Indexes[len(td.Indexes)-1]
	tableIndexes := make([]TableIndex, 0, 1)
	lastIndex.Action = &actionDrop
	tableIndexes = append(tableIndexes, lastIndex)

	lastCol := td.Cols[len(td.Cols)-1]
	tableCols := make([]TableCol, 0, 1)
	lastCol.Action = &actionDrop
	tableCols = append(tableCols, lastCol)

	td.Cols = tableCols
	td.Indexes = tableIndexes
	alterSql := td.GetAlterTableSql()
	_, _, err = store.GetFromResult(db.Exec(alterSql))
	assert.Nil(t, err)

	// add back col and index
	tableIndexes[0].Action = &actionAdd
	tableCols[0].Action = &actionAdd
	alterAddSql := td.GetAlterTableSql()
	_, _, err = store.GetFromResult(db.Exec(alterAddSql))
	assert.Nil(t, err)
	_, err = db.Exec(fmt.Sprintf("DROP TABLE %s", "user_info"))
	assert.Nil(t, err)
}

func TestTableDefinition_GenerateAlter(t *testing.T) {
	td := TableService.GetTableDefinition(db, "label", "test")
	r := td.GenerateAlter()
	assert.NotEmpty(t, r)
}
