package runner

type SqlRunner interface {
	Query() TableData
	RunDml()
	RunDDL()
}

// CellElement represents a cell element in a table
type CellElement struct {
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// TableRow represents a row in a table
type TableRow struct {
	Cells    []CellElement `json:"cells"`
	Position int           `json:"position"`
}

// TableData table data
type TableData = []TableRow

type Result struct {
	Success bool
}