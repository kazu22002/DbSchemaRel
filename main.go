package main

import (
	"./config"
	"./models"
	"./output"
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	c := config.GetConfig()

	db, err := models.OpenConn(c)
	checkError(err)
	defer db.Close()

	rows, err := models.GetSchema(db, c.GetDbName())
	checkError(err)

	output.Output(rows)
}
