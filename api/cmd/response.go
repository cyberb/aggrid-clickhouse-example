package main

type Response struct {
	Data     []Row  `json:"rowData"`
	RowCount uint64 `json:"rowCount"`
}
