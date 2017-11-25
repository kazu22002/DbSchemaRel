package schema

type DbSchema struct {
	table   string
	columns []string
}

func NewDbSchema() DbSchema {
	var ret DbSchema = DbSchema{
		table : "",
		columns : []string{""}
	}
	return ret
}
