package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//读键盘
	var (
		str    string
		reader *bufio.Reader
		err    error
	)
	fmt.Println("请输入DB路径：")
	reader = bufio.NewReader(os.Stdin)
	//以换行符结束

	str, _ = reader.ReadString('\n')

	db, err := sql.Open("sqlite3", strings.TrimSpace(str))
	if err != nil {
		fmt.Println(err, db)
	}
	fmt.Println(str + "had open for exec,please input sql to run")
	fmt.Println("input 'exit' to break the program")
	fmt.Println("please don't input 'select' statement")
	fmt.Println("support `create` `update` `insert`:")
	for {
		reader = bufio.NewReader(os.Stdin)
		str, _ = reader.ReadString('\n')
		fmt.Println("input:" + str)
		if strings.HasPrefix(str, "exit") {
			db.Close()
			fmt.Println("exit...")
			goto LOOP
		} else if strings.HasPrefix(str, "default") {
			makeDB(db)
			fmt.Println("done...")
		} else if strings.HasPrefix(str, "query ") {
			result, err := doquery(*db, strings.TrimPrefix(str, "query "))
			if err != nil {
				fmt.Println(err)
			} else {
				for key, value := range result {
					fmt.Println(key, ":", value)
				}
			}
			fmt.Println("done...")

		} else {
			_, err = db.Exec(str)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("done...")
		}
	}
LOOP:
	return

}

func makeDB(db *sql.DB) {
	_, err1 := db.Exec("CREATE TABLE `user` (`id` int PRIMARY KEY autoincrement, `name` varchar(255) NOT NULL DEFAULT '', `password` varchar(255) NOT NULL DEFAULT '', `times` int(1) NOT NULL DEFAULT '0' )")
	if err1 != nil {
		fmt.Println(err1)
	}
	_, err1 = db.Exec("CREATE TABLE `session` (`kylinseid` varchar(100) NOT NULL,`time` int(11) DEFAULT NULL,PRIMARY KEY (`kylinseid`))")
	if err1 != nil {
		fmt.Println(err1)
	}
	_, err1 = db.Exec("CREATE TABLE `notebook` (`ID` integer PRIMARY KEY autoincrement,  `theme` varchar(100) NOT NULL,`describe` varchar(1000) NOT NULL,`content` varchar(10000) DEFAULT NULL,`label` varchar(100) DEFAULT NULL,`update_time` int(11) DEFAULT NULL,`aut` varchar(1000) DEFAULT NULL)")
	if err1 != nil {
		fmt.Println(err1)
	}
}

func doquery(db sql.DB, dosql string) ([]map[string]string, error) {
	//定义一个result map数组接收真正的结果值
	var result []map[string]string
	rows, err := db.Query(dosql)
	if err != nil {
		return result, err
	}
	//定义参数接收结果
	colunms, _ := rows.Columns()
	values := make([]sql.RawBytes, len(colunms))
	scanArgs := make([]interface{}, len(colunms))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		res := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range values {
			res[colunms[i]] = string(col)
		}
		result = append(result, res)
	}
	rows.Close()
	return result, nil
}
