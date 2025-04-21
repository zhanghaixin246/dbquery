package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func isSelectQuery(query string) bool {
	q := strings.TrimLeft(strings.ToLower(query), " \t\n")
	return strings.HasPrefix(q, "select")
}

func main() {
	var (
		dsn      = flag.String("dsn", "", "Database DSN (格式: user:password@tcp(host:port)/dbname)")
		host     = flag.String("h", "localhost", "Database host")
		port     = flag.Int("P", 3306, "Database port")
		user     = flag.String("u", "", "Database user")
		password = flag.String("p", "", "Database password")
		database = flag.String("d", "", "Database name")
	)
	flag.Parse()

	connectionDSN := *dsn
	if connectionDSN == "" {
		connectionDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *password, *host, *port, *database)
	}

	db, err := sql.Open("mysql", connectionDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功！")
	fmt.Println("输入SQL查询语句，或者输入 'exit' 退出程序")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("mysql> ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		if strings.ToLower(query) == "exit" {
			fmt.Println("再见！")
			break
		}

		// 检查是否为SELECT查询
		if !isSelectQuery(query) {
			fmt.Println("错误：只允许SELECT查询操作")
			continue
		}

		// 执行查询
		rows, err := db.Query(query)
		if err != nil {
			fmt.Printf("查询错误: %v\n", err)
			continue
		}

		columns, err := rows.Columns()
		if err != nil {
			log.Printf("获取列名错误: %v\n", err)
			rows.Close()
			continue
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// 打印表头
		fmt.Println(strings.Join(columns, "\t"))
		fmt.Println(strings.Repeat("-", 80))

		// 打印数据
		rowCount := 0
		for rows.Next() {
			rowCount++
			err = rows.Scan(scanArgs...)
			if err != nil {
				log.Printf("扫描行错误: %v\n", err)
				break
			}

			var value string
			for i, col := range values {
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				if i < len(values)-1 {
					fmt.Print(value, "\t")
				} else {
					fmt.Println(value)
				}
			}
		}
		fmt.Printf("\n共 %d 行记录\n\n", rowCount)
		rows.Close()
	}
}
