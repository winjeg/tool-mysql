package worker

import (
	"github.com/stretchr/testify/assert"
	"github.com/winjeg/tool-mysql/schema"
	"testing"
)

func TestPerform(t *testing.T) {
	sqls := GenerateSql(100, 2, 10, "label", "test")
	ParallelExecute(sqls, 3)
}

func TestParallelExecute(t *testing.T) {
	tsc := schema.TableService.GetTableDefinition(db, "label", "test")
	assert.NotNil(t, tsc)
	alterSqls := tsc.GenerateAlter()
	assert.NotEmpty(t, alterSqls)
	ParallelExecute(alterSqls, 2)
}
