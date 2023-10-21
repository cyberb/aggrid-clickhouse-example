package main

import "fmt"

type FilterModel map[string]FilterDef

func (f FilterModel) ToSql() string {
	query := ""
	for field, filter := range f {
		query += fmt.Sprintf(" AND %s %s %s ", field, filter.Operator(), filter.Value())
	}

	return query
}

type FilterDef struct {
	FilterType string `json:"filterType"`
	Filter     any    `json:"filter"`
	Type       string `json:"type"`
}

func (f FilterDef) Value() string {

	switch f.FilterType {
	case "number":
		return fmt.Sprintf(`%f`, f.Filter)
	case "string":
		return fmt.Sprintf(`"%s"`, f.Filter)
	case "text":
		return fmt.Sprintf(`'%s'`, f.Filter)
	default:
		return fmt.Sprintf(`"%s"`, f.Filter)
	}
}

func (f FilterDef) Operator() string {
	switch f.Type {
	case "equals":
		return "="
	case "greaterThan":
		return ">"
	case "greaterThanOrEqual":
		return ">="
	case "LessThan":
		return "<"
	case "LessThanOrEqual":
		return "<="
	case "notEqual":
		return "!="
	default:
		return "="
	}
}
