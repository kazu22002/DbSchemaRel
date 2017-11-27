package main

import (
	"./config"
	"./dbtypes"
	"./singular"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"regexp"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func openConn(c *config.DbConf) (*sql.DB, error) {

	db, err := sql.Open("postgres", "user="+c.GetUser()+" dbname="+c.GetDbName()+" password="+c.GetPass()+" host="+c.GetHost()+" sslmode=disable")
	return db, err
}

func getSchema(db *sql.DB, db_name string) []dbtypes.DbSchema {
	table, err := getTables(db)
	checkError(err)

	length := len(table)

	var ret []dbtypes.DbSchema
	for i := 0; i < length; i++ {
		d := dbtypes.DbSchema{}

		d.SetTable(table[i])
		c, _ := getColumns(db, db_name, table[i])
		d.SetColumns(c)

		ret = append(ret, d)
	}
	return ret
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT relname as table_name FROM pg_stat_user_tables")
	checkError(err)
	defer rows.Close()

	var ret []string
	for rows.Next() {
		var table_name string
		if err := rows.Scan(&table_name); err != nil {
			checkError(err)
		}
		ret = append(ret, table_name)
	}
	return ret, err
}

func getColumns(db *sql.DB, db_name string, table_name string) ([]string, error) {

	rows, err := db.Query(fmt.Sprintf("select column_name from information_schema.columns where table_catalog = '%s' and table_name = '%s' order by ordinal_position", db_name, table_name))
	checkError(err)
	defer rows.Close()

	var ret []string
	for rows.Next() {
		var column_name string
		if err := rows.Scan(&column_name); err != nil {
			checkError(err)
		}
		if !regexp.MustCompile("id$").Match([]byte(column_name)) {
			continue
		}

		ret = append(ret, column_name)
	}
	return ret, err
}

func Output(data []dbtypes.DbSchema) {

	if len(data) < 0 {
		return
	}
	var content bytes.Buffer
	recode := "\r\n"
	count := len(data)

	content.WriteString("@startuml")
	content.WriteString(recode)
	content.WriteString("package \"データベース\" as ext <<Database>> {")
	content.WriteString(recode)
	for i := 0; i < count; i++ {
		d := data[i]
		content.WriteString("  entity \"" + d.GetTable() + "\" as " + d.GetTable() + " {")
		content.WriteString(recode)

		c := d.GetColumns()
		column_count := len(c)
		for l := 0; l < column_count; l++ {
			content.WriteString("    " + c[l])
			content.WriteString(recode)
		}
		content.WriteString("  }")
		content.WriteString(recode)
	}
	content.WriteString("}")
	content.WriteString(recode)

	for i := 0; i < count; i++ {
		d := data[i]

		single_id := singular.SingleName(d.GetTable()) + "_id"
		for l := 0; l < count; l++ {
			dd := data[l]
			cc := dd.GetColumns()
			column_count := len(cc)
			for m := 0; m < column_count; m++ {
				column_name := cc[m]
				if single_id == column_name {
					content.WriteString(d.GetTable() + " - " + dd.GetTable())
					content.WriteString(recode)
				}
			}
		}
	}

	content.WriteString("@enduml")
	content.WriteString(recode)

	ioutil.WriteFile("plant_uml.txt", []byte(content.String()), os.ModePerm)
}

func main() {
	c := config.GetConfig()

	db, err := openConn(c)
	checkError(err)
	defer db.Close()

	rows := getSchema(db, c.GetDbName())

	Output(rows)
}
