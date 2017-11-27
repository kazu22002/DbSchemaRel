package dbtypes

type DbSchema struct {
	table   string
	columns []string
}

func (d *DbSchema) SetTable(table_name string) {
	d.table = table_name
}

func (d *DbSchema) SetColumns(column_name []string) {
	d.columns = column_name
}

func (d *DbSchema) GetTable() string {
	return d.table
}

func (d *DbSchema) GetColumns() []string {
	return d.columns
}
