package parser

import (
	"fmt"
	"testing"
)

const (
	sql = `
SELECT a.name, b.title, a.name
FROM t_user as a
INNER JOIN t_comments as b ON a.user_id = b.user_id
WHERE a.user_id = (SELECT user_id FROM t_biz t WHERE c_biz=(
SELECT user_id FROM t_bizx x WHERE c_biz=1
))
`
	sql5 = `INSERT INTO user(name, value) VALUE('1', '2')`
)

func TestParser(t *testing.T) {
	v := Parse2Elements(sql)
	v2 := Parse2Elements(sql5)
	fmt.Println(v2, v)
}
