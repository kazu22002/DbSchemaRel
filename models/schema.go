package models

import (
	"../config"
	"../dbtypes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
)

func OpenConn(c *config.DbConf) (*sql.DB, error) {

	db, err := sql.Open("postgres", "user="+c.GetUser()+" dbname="+c.GetDbName()+" password="+c.GetPass()+" host="+c.GetHost()+" sslmode=disable")
	return db, err
}

func GetSchema(db *sql.DB, db_name string) ([]dbtypes.DbSchema, error) {
	table, err := getTables(db)
	if err != nil {
		return nil, err
	}

	length := len(table)

	var ret []dbtypes.DbSchema
	for i := 0; i < length; i++ {
		d := dbtypes.DbSchema{}

		d.SetTable(table[i])
		c, err := getColumns(db, db_name, table[i])
		if err != nil {
			return nil, err
		}
		d.SetColumns(c)

		ret = append(ret, d)
	}
	return ret, err
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT relname as table_name FROM pg_stat_user_tables")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []string
	for rows.Next() {
		var table_name string
		if err := rows.Scan(&table_name); err != nil {
			if err != nil {
				return nil, err
			}
		}
		ret = append(ret, table_name)
	}
	return ret, err
}

func getColumns(db *sql.DB, db_name string, table_name string) ([]string, error) {

	rows, err := db.Query(fmt.Sprintf("select column_name from information_schema.columns where table_catalog = '%s' and table_name = '%s' order by ordinal_position", db_name, table_name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []string
	for rows.Next() {
		var column_name string
		if err := rows.Scan(&column_name); err != nil {
			if err != nil {
				return nil, err
			}
		}
		if !regexp.MustCompile("id$").Match([]byte(column_name)) {
			continue
		}

		ret = append(ret, column_name)
	}
	return ret, err
}
