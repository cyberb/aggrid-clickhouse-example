package main

import "fmt"

type SortModel []SortDef

func (s SortModel) ToSql(filter map[string]bool) string {
	query := ""
	for _, order := range s {
		_, ok := filter[order.ColId]
		if len(filter) == 0 || ok {
			query += fmt.Sprintf(" %s %s ", order.Col(), order.Sort)
		}
	}
	if query != "" {
		query = " ORDER BY " + query
	}
	return query
}

type SortDef struct {
	Sort  string `json:"sort"`
	ColId string `json:"colid"`
}

func (s SortDef) Col() string {
	if s.ColId == "id" {
		return "toString(id)"
	}
	return s.ColId
}
