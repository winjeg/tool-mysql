package main

import (
	"fmt"

	"github.com/winjeg/tool-mysql/parser"
)

const sql = `
SELECT a.name, b.title, a.name
FROM t_user as a
INNER JOIN t_comments as b ON a.user_id = b.user_id
WHERE a.user_id = (SELECT user_id FROM t_biz t WHERE c_biz=(
SELECT user_id FROM t_bizx x WHERE c_biz=1
))
`

const sql2 = ` DROP TABLE a`

const sql3 = `ALTER TABLE testalter_tbl ADD i INT`

const sql4 = `UPDATE user SET naame = 'a' WHERE t = 1`

const sql5 = `INSERT INTO user(name, value) VALUE('1', '2')`

func main() {
	v := parser.Parse2Elements(sql)
	v2 := parser.Parse2Elements(sql5)
	fmt.Println(v2, v)
}
