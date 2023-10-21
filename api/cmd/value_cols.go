package main

import "fmt"

type ValueCols []ValueDef

type ValueDef struct {
	Id string `json:"id"`

	AggFunc     string `json:"aggFunc"`
	DisplayName string `json:"displayName"`
	Field       string `json:"field"`
}

func (v ValueCols) ToSql() string {
	query := ""
	for _, value := range v {
		query += fmt.Sprintf(", %s(%s) as %s", value.AggFunc, value.Field, value.Id)
	}
	return query
}
