package main

import (
	"encoding/json"
	"fmt"

	"github.com/winjeg/tool-mysql/parser"
)

const sql = `
SELECT a.name, b.title, a.name
FROM t_user as a
INNER JOIN t_comments as b ON a.user_id = b.user_id
WHERE a.user_id = (SELECT user_id FROM t_biz t WHERE c_biz=1)
`

func main() {
	astNode, err := parser.Parse(sql)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	v := parser.Extract(astNode)
	d, _ := json.Marshal(*v)
	fmt.Println(string(d))
}