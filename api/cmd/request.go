package main

import "fmt"

type Request struct {
	StartRow     int32        `json:"startRow"`
	EndRow       int32        `json:"endRow"`
	SortModel    SortModel    `json:"sortModel"`
	FilterModel  FilterModel  `json:"filterModel"`
	RowGroupCols RowGroupCols `json:"rowGroupCols"`
	ValueCols    ValueCols    `json:"valueCols"`
	GroupKeys    GroupKeys    `json:"groupKeys"`
}

func (r Request) ToGroupWhereSql() string {
	query := ""
	for _, col := range r.RowGroupCols {
		for _, key := range r.GroupKeys {
			switch key.(type) {
			case float64:
				query += fmt.Sprintf(" AND %s = %d", col.Field, int(key.(float64)))
			default:
				query += fmt.Sprintf(" AND %s = %v", col.Field, key)
			}
		}
	}

	return query
}
func (r Request) ToOrderSql() string {
	query := ""
	if r.IsGrouping() {
		query += r.SortModel.ToSql(r.RowGroupCols.Fields())
	} else {
		query += r.SortModel.ToSql(make(map[string]bool))
	}

	return query
}
func (r Request) IsGrouping() bool {
	return len(r.RowGroupCols) > len(r.GroupKeys)
}
