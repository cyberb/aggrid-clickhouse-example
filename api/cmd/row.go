package main

import "github.com/google/uuid"

type Row struct {
	Id     *uuid.UUID `json:"id" ch:"id"`
	Year   *uint16    `json:"year" ch:"year"`
	Amount *float64   `json:"amount" ch:"amount"`
	Field1 string     `json:"field1" ch:"field1"`
	Field2 string     `json:"field2" ch:"field2"`
	Field3 string     `json:"field3" ch:"field3"`
	Field4 string     `json:"field4" ch:"field4"`
	Field5 string     `json:"field5" ch:"field5"`
	Field6 string     `json:"field6" ch:"field6"`
}
