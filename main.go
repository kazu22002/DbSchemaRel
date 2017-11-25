package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"regexp"
)

type DbSchema struct {
	table   string
	columns []string
}

type DbConf struct {
	User   string `json:"user"`
	DbName string `json:"db_name"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func getConfig() (DbConf, error) {
	configFile, err := os.Open("config.json")
	checkError(err)
	decoder := json.NewDecoder(configFile)
	var config DbConf
	err = decoder.Decode(&config)
	checkError(err)
	defer configFile.Close()

	return config, err
}

func openConn(config DbConf) (*sql.DB, error) {

	db, err := sql.Open("postgres", "user="+config.User+" dbname="+config.DbName+" password="+config.Pass+" host="+config.Host+" sslmode=disable")
	return db, err
}

func getSchema(db *sql.DB, db_name string) []DbSchema {
	table, err := getTables(db)
	checkError(err)

	length := len(table)

	var ret []DbSchema
	for i := 0; i < length; i++ {
		var d DbSchema

		d.table = table[i]
		d.columns, err = getColumns(db, db_name, table[i])

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

var singular_rules = map[string]string{
	"(s)tatuses$":     "12tatus",
	"^(.*)(menu)s$":   "12",
	"(quiz)zes$":      "1",
	"(matr)ices$":     "1ix",
	"(vert|ind)ices$": "1ex",
	"^(ox)en":         "1",
	"(alias)(es)*$":   "1",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$": "1us",
	"([ftw]ax)es":        "1",
	"(cris|ax|test)es$":  "1is",
	"(shoe|slave)s$":     "1",
	"(o)es$":             "1",
	"ouses$":             "ouse",
	"([^a])uses$":        "1us",
	"([m|l])ice$":        "1ouse",
	"(x|ch|ss|sh)es$":    "1",
	"(m)ovies$":          "12ovie",
	"(s)eries$":          "12eries",
	"([^aeiouy]|qu)ies$": "1y",
	"([lr])ves$":         "1f",
	"(tive)s$":           "1",
	"(hive)s$":           "1",
	"(drive)s$":          "1",
	"([^fo])ves$":        "1fe",
	"(^analy)ses$":       "1sis",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$": "12sis",
	"([ti])a$":    "1um",
	"(p)eople$":   "12erson",
	"(m)en$":      "1an",
	"(c)hildren$": "12hild",
	"(n)ews$":     "12ews",
	"eaus$":       "eau",
	"^(.*us)$":    "1",
	"s$":          ""}

func singleName(name string) string {
	var single_name = name

	for key, replace := range singular_rules {
		if regexp.MustCompile(key).MatchString(name) {
			single_name = regexp.MustCompile(key).ReplaceAllString(name, replace)
			break
		}
	}

	return single_name
}

func Output(data []DbSchema) {

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
		content.WriteString("  entity \"" + data[i].table + "\" " + data[i].table + " {")
		content.WriteString(recode)

		column_count := len(data[i].columns)
		for l := 0; l < column_count; l++ {
			content.WriteString("    " + data[i].columns[l])
			content.WriteString(recode)
		}
		content.WriteString("  }")
		content.WriteString(recode)
	}
	content.WriteString("}")
	content.WriteString(recode)
	content.WriteString("@enduml")
	content.WriteString(recode)

	ioutil.WriteFile("plant_uml.txt", []byte(content.String()), os.ModePerm)

	// single_name := singleName(name)
}

func main() {
	config, err := getConfig()

	db, err := openConn(config)
	checkError(err)
	defer db.Close()

	rows := getSchema(db, config.DbName)

	Output(rows)
}
