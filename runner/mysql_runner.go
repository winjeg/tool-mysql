package runner

import (
	"database/sql"
	"strconv"
	"strings"
)

// QueryData Query data from a db
func QueryData(db *sql.DB, query string) (TableData, error) {
	rows, qErr := db.Query(query)
	if qErr != nil {
		return nil, qErr
	}
	cols, rErr := rows.Columns()
	if rErr != nil {
		return nil, rErr
	}
	colNum := len(cols)
	colTypes, ctErr := rows.ColumnTypes()
	if ctErr != nil {
		return nil, ctErr
	}
	values := make([]*[]byte, colNum)
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	result := make([]TableRow, 0, 100)
	currentRow := 0
	for rows.Next() {
		if scanErr := rows.Scan(scans...); scanErr != nil {
			return nil, scanErr
		}
		var row TableRow
		row.Position = currentRow
		cells := make([]CellElement, colNum)
		for i, v := range values {
			cells[i].Key = cols[i]
			cells[i].Type = strings.ToUpper(colTypes[i].DatabaseTypeName())
			if v == nil {
				cells[i].Value = nil
			} else {
				switch cells[i].Type {
				case "INT", "TINYINT", "BIGINT", "MEDIUMINT", "SMALLINT":
					cells[i].Value, _ = strconv.Atoi(string(*v))
				case "DECIMAL", "FLOAT", "DOUBLE":
					cells[i].Value, _ = strconv.ParseFloat(string(*v), 64)
				default:
					cells[i].Value = string(*v)
				}
			}
		}
		row.Cells = cells
		result = append(result, row)
		currentRow++
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)
	return result, nil
}
