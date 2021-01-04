package parser

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"

)

type colX struct{
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
func Extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}
func Parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes[0], nil
}
