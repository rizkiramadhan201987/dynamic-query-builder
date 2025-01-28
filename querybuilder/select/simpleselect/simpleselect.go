package simpleselect

import (
	"fmt"
	"strings"
)

// SimpleSelect implements basic SELECT clause
type SimpleSelect struct {
	fields []string
}

func NewSimpleSelect(fields ...string) *SimpleSelect {
	if len(fields) == 0 {
		fields = []string{"*"}
	}
	return &SimpleSelect{
		fields: fields,
	}
}

func (s *SimpleSelect) Build() (string, error) {
	if len(s.fields) == 0 {
		return "", fmt.Errorf("no fields specified for select")
	}
	return fmt.Sprintf("SELECT %s", strings.Join(s.fields, ", ")), nil
}
