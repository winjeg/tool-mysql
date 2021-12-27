package parser

import (
	"fmt"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
)

type SelectField struct {
	TableName string `json:"tableName"`
	ColName   string `json:"colName"`
	AsName    string `json:"asName"`
}

type SqlType int

const (
	ReadStmt  = SqlType(1)
	DmlStmt   = SqlType(2)
	DdlStmt   = SqlType(3)
	OtherStmt = SqlType(9)
)

var (
	tiParser = parser.New()
)

type sqlElements struct {
	SelectFields []SelectField     `json:"selectFields"`
	TableNames   map[string]string `json:"tableNames"`
	Type         SqlType           `json:"type"`
	StmtNum      int               `json:"stmtNum"`
}

func (v *sqlElements) getCols(in ast.Node) {
	if selectField, ok := in.(*ast.SelectField); ok {
		if v.SelectFields == nil {
			v.SelectFields = make([]SelectField, 0, 10)
		}
		s := selectField.Expr
		if se, ok := s.(*ast.ColumnNameExpr); ok {
			colName := se.Name
			sf := SelectField{
				TableName: colName.Table.String(),
				ColName:   colName.Name.String(),
				AsName:    selectField.AsName.String(),
			}
			v.SelectFields = append(v.SelectFields, sf)
		}
	}
}

func (v *sqlElements) getTables(in ast.Node) {
	if table, ok := in.(*ast.TableSource); ok {
		if v.TableNames == nil {
			v.TableNames = make(map[string]string, 10)
		}
		source := table.Source
		if tableName, ok := source.(*ast.TableName); ok {
			v.TableNames[tableName.Name.String()] = table.AsName.String()
		}
	}
}

func (v *sqlElements) getSubQueries(in ast.Node) {
	if _, ok := in.(*ast.SelectStmt); ok {
		v.StmtNum++
	}
}

func (v *sqlElements) Enter(in ast.Node) (ast.Node, bool) {
	v.getCols(in)
	v.getTables(in)
	v.getSubQueries(in)
	return in, false
}

func (v *sqlElements) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func Extract(rootNode *ast.StmtNode) *sqlElements {
	if rootNode == nil {
		return nil
	}
	v := &sqlElements{}
	switch t := (*rootNode).(type) {
	case *ast.SelectStmt:
		v.Type = ReadStmt
	case *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
		v.Type = DmlStmt
	case *ast.CreateIndexStmt, *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.AlterTableStmt, *ast.DropTableStmt, *ast.DropIndexStmt, *ast.AlterDatabaseStmt, *ast.TruncateTableStmt:
		v.Type = DdlStmt
	default:
		_ = t
		v.Type = OtherStmt
	}
	(*rootNode).Accept(v)
	return v
}

func Parse(sql string) (*ast.StmtNode, error) {
	stmtNodes, _, err := tiParser.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes[0], nil
}

func Parse2Elements(sql string) *sqlElements {
	astNode, err := Parse(sql)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil
	}
	return Extract(astNode)
}
