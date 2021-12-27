package schema

import (
	"fmt"
	"github.com/winjeg/tool-mysql/generator"
	"strings"
)

const (
	maxStrLen = 255
	// data types
	intUnsigned    = "UNSIGNED"
	typeInt        = "INT"
	typeTinyInt    = "TINYINT"
	typeSmallInt   = "SMALLINT"
	typeMediumInt  = "MEDIUMINT"
	typeBigInt     = "BIGINT"
	typeFloat      = "FLOAT"
	typeDouble     = "DOUBLE"
	typeDecimal    = "DECIMAL"
	typeVarchar    = "VARCHAR"
	typeChar       = "CHAR"
	typeTinyText   = "TINYTEXT"
	typeMediumText = "MEDIUMTEXT"
	typeText       = "TEXT"
	typeLongText   = "LONGTEXT"
	typeDateTime   = "DATETIME"
	typeTimestamp  = "TIMESTAMP"
	typeBit        = "BIT"
)

// GetValue get value for constructing insert fragment
func (t *TableCol) GetValue() string {
	dataType := strings.ToUpper(t.DataType)
	if t.IsUnsigned() {
		// only numbers can be unsigned, i bet
		switch dataType {
		case typeInt:
			return generator.IntRander.Get(generator.TypeUInt)
		case typeTinyInt:
			return generator.IntRander.Get(generator.TypeUTinyInt)
		case typeSmallInt:
			return generator.IntRander.Get(generator.TypeSmallInt)
		case typeMediumInt:
			return generator.IntRander.Get(generator.TypeMediumInt)
		case typeBigInt:
			return generator.IntRander.Get(generator.TypeBigInt)
		case typeDouble, typeFloat, typeDecimal:
			return generator.FloatRander.Get(int(*t.NumPrecision)-int(*t.NumScale), int(*t.NumScale))
		default:
			return generator.IntRander.Get(generator.TypeUTinyInt)
		}
	} else {
		// all the other types may not unsigned
		switch dataType {
		case typeInt:
			return generator.IntRander.Get(generator.TypeInt)
		case typeTinyInt:
			return generator.IntRander.Get(generator.TypeTinyInt)
		case typeSmallInt:
			return generator.IntRander.Get(generator.TypeSmallInt)
		case typeMediumInt:
			return generator.IntRander.Get(generator.TypeMediumInt)
		case typeBigInt:
			return generator.IntRander.Get(generator.TypeBigInt)
		case typeDouble:
			return randomFloatStr(generator.FloatRander.GetFloat(int(*t.NumPrecision)-int(*t.NumScale), int(*t.NumScale)))
		case typeFloat:
			return randomFloatStr(generator.FloatRander.GetFloat(int(*t.NumPrecision)-int(*t.NumScale), int(*t.NumScale)))
		case typeDecimal:
			return randomFloatStr(generator.FloatRander.GetFloat(int(*t.NumPrecision)-int(*t.NumScale), int(*t.NumScale)))
		case typeText, typeLongText, typeTinyText, typeMediumText, typeChar, typeVarchar:
			return randomStr(uint(*t.CharMaxLen))
		case typeTimestamp, typeDateTime:
			return generator.TimeRander.Get()
		case typeBit:
			return generator.BitRander.Get(int(*t.NumPrecision))

		default:
			return generator.IntRander.Get(generator.TypeUTinyInt)
		}
	}
}

// IsUnsigned to judge if an column is unsigned or not
func (t *TableCol) IsUnsigned() bool {
	return strings.Index(strings.ToUpper(t.ColType), intUnsigned) > 0
}

// to convert a float string to its
func randomFloatStr(s string) string {
	if generator.RandomBool() {
		return fmt.Sprintf("'-%s'", s)
	}
	return fmt.Sprintf("'%s'", s)
}

// for the string generation may be too slow if the length of the text is too long
// we here by make the generated string's length is at most 255
func randomStr(maxLen uint) string {
	if maxLen > maxStrLen {
		maxLen = maxStrLen
	}
	n := generator.RandomUInt(uint64(maxLen))
	return generator.StringRander.Get(int(n))
}

func (td *TableDefinition) GenerateInsert(n int) string {
	if len(td.Cols) < 1 {
		return ""
	}
	if n < 1 {
		n = 1
	}
	return td.genInsertDef() + "\n" + td.genValues(n) + ";"
}

func (td *TableDefinition) genInsertDef() string {
	if len(td.Cols) < 1 {
		return ""
	}
	var colDef string
	for i := range td.Cols {
		element := fmt.Sprintf("`%s`,", td.Cols[i].ColName)
		colDef += element
	}
	return fmt.Sprintf("INSERT INTO `%s`(%s) VALUES", td.TableName, colDef[:len(colDef)-1])
}

func (td *TableDefinition) genValues(n int) string {
	if len(td.Cols) < 1 {
		return ""
	}
	if n < 1 {
		n = 1
	}
	var lines string
	for x := 0; x < n; x++ {
		var colVal string
		for i := range td.Cols {
			colVal += fmt.Sprintf("%s,", td.Cols[i].GetValue())
		}
		lines += fmt.Sprintf("(%s),\n", colVal[:len(colVal)-1])
	}
	return lines[:len(lines)-2]
}

// GenerateAlter generate alter table sql
// currently the generated sql is symmetric, when all executed successfully,
// the table structure should be in its original form
func (td *TableDefinition) GenerateAlter() []string {
	indexLen := len(td.Indexes)
	colLen := len(td.Cols)
	result := make([]string, 0, defaultSize)
	tableCols := make([]TableCol, 0, 1)
	tableIndexes := make([]TableIndex, 0, 1)
	oriCols := td.Cols
	oriIndexes := td.Indexes
	td.Extra.TableComment = ""

	if indexLen > 0 {
		lastIndex := oriIndexes[len(oriIndexes)-1]
		lastIndex.Action = &actionDrop
		tableIndexes = append(tableIndexes, lastIndex)
		td.Cols = nil
		td.Indexes = tableIndexes
		sql1 := td.GetAlterTableSql() + ";"
		td.Indexes[0].Action = &actionAdd
		td.Cols = nil
		td.Indexes = tableIndexes
		sql2 := td.GetAlterTableSql() + "';"

		result = append(result, sql1, sql2)
	}
	// generate two  col relate alter
	if colLen > 0 {
		lastCol := oriCols[len(oriCols)-1]
		lastCol.Action = &actionDrop
		tableCols = append(tableCols, lastCol)
		td.Cols = tableCols
		td.Indexes = nil
		sql1 := td.GetAlterTableSql() + ";"
		tableCols[0].Action = &actionAdd
		td.Cols = tableCols
		td.Indexes = nil
		sql2 := td.GetAlterTableSql() + ";"
		result = append(result, sql1, sql2)
	}
	return result
}
