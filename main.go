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

/**
 * ループで回した時に、ランダムに抽出されたためキーのみ配列を作成する。
 */
var singular_rules = map[string]string{
	"(s)tatuses$":     "${1}tatus",
	"^(.*)(menu)s$":   "${1}",
	"(quiz)zes$":      "${1}",
	"(matr)ices$":     "${1}ix",
	"(vert|ind)ices$": "${1}ex",
	"^(ox)en":         "${1}",
	"(alias)(es)*$":   "${1}",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$": "${1}us",
	"([ftw]ax)es":        "${1}",
	"(cris|ax|test)es$":  "${1}is",
	"(shoe|slave)s$":     "${1}",
	"(o)es$":             "${1}",
	"ouses$":             "ouse",
	"([^a])uses$":        "${1}us",
	"([m|l])ice$":        "${1}ouse",
	"(x|ch|ss|sh)es$":    "${1}",
	"(m)ovies$":          "${1}ovie",
	"(s)eries$":          "${1}eries",
	"([^aeiouy]|qu)ies$": "${1}y",
	"([lr])ves$":         "${1}f",
	"(tive)s$":           "${1}",
	"(hive)s$":           "${1}",
	"(drive)s$":          "${1}",
	"([^fo])ves$":        "${1}fe",
	"(^analy)ses$":       "${1}sis",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$": "${1}sis",
	"([ti])a$":    "${1}um",
	"(p)eople$":   "${1}erson",
	"(m)en$":      "${1}an",
	"(c)hildren$": "${1}hild",
	"(n)ews$":     "${1}ews",
	"eaus$":       "eau",
	"^(.*us)$":    "${1}",
	"s$":          ""}

var singular_rules_sort = []string{
	"(s)tatuses$",
	"^(.*)(menu)s$",
	"(quiz)zes$",
	"(matr)ices$",
	"(vert|ind)ices$",
	"^(ox)en",
	"(alias)(es)*$",
	"(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$",
	"([ftw]ax)es",
	"(cris|ax|test)es$",
	"(shoe|slave)s$",
	"(o)es$",
	"ouses$",
	"([^a])uses$",
	"([m|l])ice$",
	"(x|ch|ss|sh)es$",
	"(m)ovies$",
	"(s)eries$",
	"([^aeiouy]|qu)ies$",
	"([lr])ves$",
	"(tive)s$",
	"(hive)s$",
	"(drive)s$",
	"([^fo])ves$",
	"(^analy)ses$",
	"(analy|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$",
	"([ti])a$",
	"(p)eople$",
	"(m)en$",
	"(c)hildren$",
	"(n)ews$",
	"eaus$",
	"^(.*us)$",
	"s$",
}

func singleName(name string) string {
	var single_name = name

	for _, key := range singular_rules_sort {
		if regexp.MustCompile(key).MatchString(name) {
			single_name = regexp.MustCompile(key).ReplaceAllString(name, singular_rules[key])
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
		content.WriteString("  entity \"" + data[i].table + "\" as " + data[i].table + " {")
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

	for i := 0; i < count; i++ {
		single_id := singleName(data[i].table) + "_id"
		for l := 0; l < count; l++ {
			column_count := len(data[l].columns)
			for m := 0; m < column_count; m++ {
				column_name := data[l].columns[m]
				if single_id == column_name {
					content.WriteString(data[i].table + " - " + data[l].table)
					content.WriteString(recode)
				}
			}
		}
	}

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
