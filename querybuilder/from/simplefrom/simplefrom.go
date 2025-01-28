package simplefrom

import "fmt"

// SimpleFrom implements a basic FROM clause
type SimpleFrom struct {
	Table string
}

func NewSimpleFrom(table string) *SimpleFrom {
	return &SimpleFrom{Table: table}
}

func (f *SimpleFrom) Build() string {
	return fmt.Sprintf("FROM %s", f.Table)
}
