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
	READ_STMT = SqlType(1)
	DML_STMT  = SqlType(2)
	DDL_STMT  = SqlType(3)
)

type sqlElements struct {
	SelectFields []SelectField     `json:"selectFields"`
	TableNames   map[string]string `json:"tableNames"`
	Type         SqlType
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
	if table, ok := in.(*ast.SelectStmt); ok {
		// -> 作为一个工具
		fmt.Println(table.Text())
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
	v := &sqlElements{}
	(*rootNode).Accept(v)
	return v
}
func Parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes[0], nil
}
