package schema

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	// default size for slices
	defaultSize = 10
)

var (
	// actions
	actionDrop   = "DROP"
	actionAdd    = "ADD"
	actionChange = "CHANGE"
	actionModify = "MODIFY"

	colTemplate = "%s COLUMN"
)

type schemaService struct {
}

var TableService = schemaService{}

const (
	tableExtra   = "SELECT AUTO_INCREMENT, TABLE_COMMENT  FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s';"
	tableIndexes = `SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME, SEQ_IN_INDEX,  INDEX_TYPE, INDEX_COMMENT
					FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'`
	tableCols = `SELECT COLUMN_NAME, ORDINAL_POSITION,  COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE,  
					CHARACTER_MAXIMUM_LENGTH,  NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT
					FROM INFORMATION_SCHEMA.COLUMNS
					WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'`
)

// TableExtra table extra information
type TableExtra struct {
	Action        *string `json:"action,omitempty"` // nil (unmodified), change (modify)
	AutoIncrement *uint32 `json:"autoIncrement,omitempty"`
	TableComment  string  `json:"tableComment"`
}

// TableIndex table index structure
type TableIndex struct {
	Action       *string `json:"action,omitempty"` // the action to add/change/drop index
	IndexName    string  `json:"indexName"`
	IndexComment string  `json:"indexComment"`
	IndexType    string  `json:"indexType"`
	ColName      string  `json:"colName"`
	ColSeq       uint32  `json:"colSeq"`
	NonUnique    bool    `json:"nonUnique"`
}

// for sorting the table index
type tbIdxes []*TableIndex

func (I tbIdxes) Len() int {
	return len(I)
}
func (I tbIdxes) Less(i, j int) bool {
	return I[i].ColSeq < I[j].ColSeq
}
func (I tbIdxes) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

// TableCol schema column information from INFORMATION_SCHEMA.COLUMNS
type TableCol struct {
	Action       *string `json:"action,omitempty"` //表结构变更类型，nil（未变），ADD（new），CHANGE（modify）
	Position     *string `json:"position,omitempty"`
	ColName      string  `json:"colName"`
	OldColName   *string `json:"oldColName,omitempty"`
	ColPosition  uint32  `json:"colPosition"`
	ColDefault   *string `json:"colDefault,omitempty"`
	IsNullable   string  `json:"isNullable"`
	DataType     string  `json:"dataType"`
	CharMaxLen   *uint32 `json:"charMaxLen,omitempty"`
	NumPrecision *uint32 `json:"numPrecision,omitempty"`
	NumScale     *uint32 `json:"numScale,omitempty"`
	ColType      string  `json:"colType"`
	ColKey       string  `json:"colKey"`
	Extra        string  `json:"extra"`
	ColComment   string  `json:"colComment"`
}

type TableDefinition struct {
	TableName string       `json:"tableName"`
	Cols      []TableCol   `json:"cols"`
	Indexes   []TableIndex `json:"indexes"`
	Extra     TableExtra   `json:"extra"`
}

// GetTableDefinition get whole table definition from database
func (ss *schemaService) GetTableDefinition(myDb *sql.DB, tableName, dbName string) TableDefinition {
	return TableDefinition{
		TableName: tableName,
		Extra:     ss.GetTableExtra(myDb, tableName, dbName),
		Cols:      ss.GetTableCols(myDb, tableName, dbName),
		Indexes:   ss.GetTableIndexes(myDb, tableName, dbName),
	}
}

// GetTableExtra get table extra info
func (*schemaService) GetTableExtra(myDb *sql.DB, tableName, dbName string) TableExtra {
	rows := query(myDb, tableName, dbName, tableExtra)
	defer rows.Close()
	extras := make([]TableExtra, 0, defaultSize)
	for rows.Next() {
		var extra TableExtra
		err := rows.Scan(&extra.AutoIncrement, &extra.TableComment)
		if err != nil {
			continue
		}
		extras = append(extras, extra)
	}
	if len(extras) > 0 {
		return extras[0]
	} else {
		var extra TableExtra
		return extra
	}
}

// GetTableCols get table cols from database
func (*schemaService) GetTableCols(myDb *sql.DB, tableName, dbName string) []TableCol {
	if myDb == nil || len(tableName) < 1 || len(dbName) < 1 {
		return nil
	}
	rows := query(myDb, tableName, dbName, tableCols)
	defer rows.Close()
	cols := make([]TableCol, 0, defaultSize)
	for rows.Next() {
		var col TableCol
		err := rows.Scan(&col.ColName, &col.ColPosition, &col.ColDefault, &col.IsNullable, &col.DataType, &col.CharMaxLen,
			&col.NumPrecision, &col.NumScale, &col.ColType, &col.ColKey, &col.Extra, &col.ColComment)
		if err != nil {
			continue
		}
		cols = append(cols, col)
	}
	return cols
}

// GetTableIndexes get table indexes from database
func (*schemaService) GetTableIndexes(myDb *sql.DB, tableName, dbName string) []TableIndex {
	rows := query(myDb, tableName, dbName, tableIndexes)
	defer rows.Close()
	indexes := make([]TableIndex, 0, defaultSize)
	for rows.Next() {
		var index TableIndex
		err := rows.Scan(&index.IndexName, &index.NonUnique, &index.ColName, &index.ColSeq, &index.IndexType,
			&index.IndexComment)
		if err != nil {
			continue
		}
		indexes = append(indexes, index)
	}
	return indexes
}

// query data from db
func query(myDb *sql.DB, tableName, dbName, sqlTemplate string) *sql.Rows {
	q := fmt.Sprintf(sqlTemplate, dbName, tableName)
	rows, err := myDb.Query(q)
	if err != nil {
		log.Printf("%v", err)
	}
	return rows
}

// GetCreateTableSql construct create table ddl from table definition
func (td *TableDefinition) GetCreateTableSql() string {
	if len(td.TableName) < 1 {
		return ""
	}
	result := fmt.Sprintf("CREATE TABLE `%s`(\n", td.TableName)
	cols := td.Cols
	for _, v := range cols {
		result += constructTableColSqlFromColDefinition(v) + ",\n"
	}
	if len(td.Indexes) < 1 {
		r := strings.TrimSpace(result)
		if strings.LastIndex(r, ",") == len(r)-1 {
			result = r[0 : len(r)-1]
		}
	}
	result += constructTableIndexFromTableIndex(td.Indexes)
	if strings.Index(result, ",") == len(result)-1 {
		result = result[:len(result)-1]
	}
	var autoIncSql string
	if td.Extra.AutoIncrement != nil {
		autoIncSql = fmt.Sprintf("AUTO_INCREMENT=%d ", *td.Extra.AutoIncrement)
	}
	extra := fmt.Sprintf(") ENGINE=InnoDB %s DEFAULT CHARSET=utf8 COMMENT='%s'",
		autoIncSql, td.Extra.TableComment)
	return result + extra
}

// construct table col sql fragments from col definition
func constructTableColSqlFromColDefinition(col TableCol) string {
	var colPrefix string
	if col.Action != nil {
		switch strings.ToUpper(*col.Action) {
		case actionDrop:
			colPrefix = fmt.Sprintf(colTemplate, actionDrop)
		case actionChange:
			colPrefix = fmt.Sprintf(colTemplate, actionChange)
		case actionAdd:
			colPrefix = fmt.Sprintf(colTemplate, actionAdd)
		case actionModify:
			colPrefix = fmt.Sprintf(colTemplate, actionModify)
		default:
			colPrefix = ""
		}
	}

	var colSuffix string
	if col.Position != nil {
		colSuffix = *col.Position
	}
	result := colPrefix
	if strings.EqualFold(colPrefix, "DROP COLUMN") {
		return fmt.Sprintf("DROP COLUMN `%s`", col.ColName)
	}

	if strings.EqualFold(colPrefix, "CHANGE COLUMN") && col.OldColName != nil {
		result += " `" + *col.OldColName + "`" + " `" + col.ColName + "` "
	} else {
		result += " `" + col.ColName + "`" + " "
	}

	result += col.ColType + " "
	if strings.EqualFold(col.IsNullable, "NO") {
		result += " NOT NULL "
	} else {
		result += " NULL "
	}

	// dealing with default values of cols
	if col.ColDefault != nil {
		if !strings.EqualFold(strings.ToUpper(*col.ColDefault), "NULL") {
			colType := strings.ToLower(col.ColType)
			// if it is of type char or type text, the default value should be surrounded by quotes
			if strings.Contains(colType, "char") || strings.Contains(colType, "text") {
				defaultVal := strings.Replace(*col.ColDefault, "\"", "\\\"", -1)
				defaultVal = strings.Replace(defaultVal, "'", "\\'", -1)
				result += "DEFAULT " + "'" + defaultVal + "' "
			} else {
				result += "DEFAULT " + *col.ColDefault + " "
			}
		} else if !strings.EqualFold(col.IsNullable, "NO") {
			result += "DEFAULT NULL "
		}
	}

	// dealing with extra information and column comments
	result += strings.ToUpper(col.Extra)
	if len(strings.TrimSpace(col.ColComment)) > 0 {
		result += fmt.Sprintf(" COMMENT '%s'", col.ColComment)
	}
	return result + colSuffix
}

// construct sql fragments for index part from table index slice
func constructTableIndexFromTableIndex(indexes []TableIndex) string {
	if len(indexes) < 1 {
		return ""
	}
	var result string
	// put all index into a map indexed by index name
	// because of some index may need more than one column, so the map value is a slice
	keyMap := make(map[string][]*TableIndex, 10)
	for _, v := range indexes {
		if i, ok := keyMap[v.IndexName]; ok {
			x := TableIndex(v)
			i = append(i, &x)
			keyMap[v.IndexName] = i
		} else {
			val := make([]*TableIndex, 0, 10)
			d := TableIndex(v)
			val = append(val, &d)
			keyMap[v.IndexName] = val
		}
	}

	for k, v := range keyMap {
		var idxPrefix string
		// sort to make sure the combination index get the right order
		sort.Sort(tbIdxes(v))
		// dealing with primary key
		if strings.EqualFold(k, "PRIMARY") {
			pkTmpl := "%s PRIMARY KEY (%s),"
			keyVal := ""
			for i, n := range v {
				if n.Action != nil && strings.EqualFold(*n.Action, actionDrop) {
					result += "DROP PRIMARY KEY,\n"
					idxPrefix = actionDrop
					continue
				}
				if n.Action != nil && strings.EqualFold(*n.Action, actionAdd) {
					idxPrefix = actionAdd
				}
				keyVal += fmt.Sprintf("`%s`", n.ColName)
				if i < len(v)-1 {
					keyVal += ","
				}
			}
			if len(keyVal) > 0 {
				result += fmt.Sprintf(pkTmpl, idxPrefix, keyVal)
			}
		} else {
			// dealing with other type of indexes
			var keyPrefix, keyVal string
			for i, n := range v {
				if n.Action != nil && strings.EqualFold(*n.Action, actionDrop) {
					result += fmt.Sprintf("DROP INDEX `%s`,", k)
					idxPrefix = actionDrop
					continue
				}
				// if it is a unique index
				if !n.NonUnique {
					keyPrefix = "UNIQUE"
				} else if strings.EqualFold(v[0].IndexType, "FULLTEXT") {
					keyPrefix = "FULLTEXT"
				}
				if n.Action != nil && strings.EqualFold(*n.Action, actionAdd) {
					idxPrefix = actionAdd
				}
				keyVal += fmt.Sprintf("`%s`", n.ColName)
				if i < len(v)-1 {
					keyVal += ","
				}
			}

			// if not for dropping index, join the sql fragments
			if !strings.EqualFold(idxPrefix, actionDrop) {
				c := fmt.Sprintf("%s %s KEY `%s` (%s)", idxPrefix, keyPrefix, k, keyVal)
				result += c
				if strings.EqualFold(strings.ToUpper(v[0].IndexType), "BTREE") {
					result += " USING BTREE "
				}
				if len(strings.TrimSpace(v[0].IndexComment)) > 0 {
					result += fmt.Sprintf(" COMMENT '%s'", v[0].IndexComment)
				}
				result += ","
			}
		}
	}
	result = strings.TrimSpace(result)
	// make sure no extra comma
	if strings.LastIndex(result, ",") == len(result)-1 {
		return result[0 : len(result)-1]
	}
	return result
}

// GetAlterTableSql get alter table sql from table definition
func (td *TableDefinition) GetAlterTableSql() string {
	if td.Cols == nil && td.Indexes == nil && td.Extra.Action == nil {
		return ""
	}

	// header
	result := fmt.Sprintf("ALTER TABLE `%s`\n", td.TableName)

	// the columns
	if len(td.Cols) > 0 {
		for _, v := range td.Cols {
			result += constructTableColSqlFromColDefinition(v) + ",\n"
		}
	}
	// indexes
	if len(td.Indexes) > 0 {
		result += constructTableIndexFromTableIndex(td.Indexes) + ","
	}
	// table comments
	if len(td.Extra.TableComment) > 0 {
		result += fmt.Sprintf("COMMENT='%s',", td.Extra.TableComment)
	}
	// auto increment
	if td.Extra.AutoIncrement != nil && *td.Extra.AutoIncrement > 0 {
		result += fmt.Sprintf("AUTO_INCREMENT=%d;", *td.Extra.AutoIncrement)
	}
	result = strings.TrimSpace(result)
	if strings.LastIndex(result, ",") == len(result)-1 {
		return result[0 : len(result)-1]
	}
	return result
}
