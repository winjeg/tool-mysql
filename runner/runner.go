package runner

type SqlRunner interface {
	Query() TableData
	RunDml()
	RunDDL()
}

// represents a cell element in a table
type CellElement struct {
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// represents a row in a table
type TableRow struct {
	Cells    []CellElement `json:"cells"`
	Position int           `json:"position"`
}

// table data
type TableData = []TableRow

type Result struct {
	Success bool
}