package main

import (
	"encoding/json"
	"fmt"

	"github.com/winjeg/tool-mysql/runner"
	"github.com/winjeg/tool-mysql/store"
)

func main() {
	db := store.GetDb()
	data, _ := runner.QueryData(db, "SELECT * FROM user_info")
	jd, _ := json.Marshal(data)
	fmt.Println(string(jd))
}
