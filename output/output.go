package output

import (
	"../dbtypes"
	"../singular"
	"bytes"
	"io/ioutil"
	"os"
)

func Output(data []dbtypes.DbSchema) {

	if len(data) < 0 {
		return
	}

	var content bytes.Buffer

	content.WriteString(outputHead())
	content.WriteString(outputEntity(data))
	content.WriteString(outputRelational(data))
	content.WriteString(outputFooter())

	ioutil.WriteFile("plant_uml.txt", []byte(content.String()), os.ModePerm)
}

func outputHead() string {
	var content bytes.Buffer
	content.WriteString("@startuml")
	content.WriteString("\r\n")
	return content.String()
}
func outputFooter() string {
	var content bytes.Buffer
	content.WriteString("@enduml")
	content.WriteString("\r\n")
	return content.String()
}

func outputEntity(data []dbtypes.DbSchema) string {
	var content bytes.Buffer
	recode := "\r\n"
	count := len(data)

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

	return content.String()
}

func outputRelational(data []dbtypes.DbSchema) string {
	var content bytes.Buffer
	recode := "\r\n"
	count := len(data)

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
	return content.String()
}
