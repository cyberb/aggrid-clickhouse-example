package main

import "fmt"

type RowGroupCols []RowGroupDef

type RowGroupDef struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Field       string `json:"field"`
}

func (r RowGroupCols) IsEmpty() bool {
	return len(r) == 0
}

func (r RowGroupCols) ToSelectSql() string {

	query := ""
	separator := " "
	for _, value := range r {
		query += fmt.Sprintf("%s %s", separator, value.Field)
		separator = ", "
	}
	return query
}
func (r RowGroupCols) ToGroupSql() string {
	query := ""
	if !r.IsEmpty() {
		query = " GROUP BY "
	}
	for _, value := range r {
		query += fmt.Sprintf(" %s ", value.Field)
	}
	return query
}

func (r RowGroupCols) Fields() map[string]bool {
	fields := make(map[string]bool)
	for _, group := range r {
		fields[group.Field] = true
	}

	return fields
}
