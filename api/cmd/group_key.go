package main

type GroupKeys []any

func (k GroupKeys) IsEmpty() bool {
	return len(k) == 0
}
