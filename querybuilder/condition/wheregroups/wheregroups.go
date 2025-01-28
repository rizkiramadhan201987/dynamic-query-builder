package wheregroups

import (
	"dynamic-sqlbuilder/querybuilder"
	"fmt"
	"strings"
)

// WhereGroups implements WhereClause
type WhereGroups struct {
	Groups []querybuilder.WhereGroup
}

func NewWhereGroups() *WhereGroups {
	return &WhereGroups{
		Groups: make([]querybuilder.WhereGroup, 0),
	}
}
func (w *WhereGroups) Add(group querybuilder.WhereGroup) {
	// fmt.Printf("Adding group with %d conditions to WhereGroups\n", len(group.Conditions))
	w.Groups = append(w.Groups, group)
}

func (w *WhereGroups) Build(paramOffset int) (string, []interface{}, error) {
	if len(w.Groups) == 0 {
		return "", nil, nil
	}

	var clauses []string
	var args []interface{}
	currentOffset := paramOffset

	for i, group := range w.Groups {
		clause, groupArgs, err := group.Build(currentOffset)
		if err != nil {
			return "", nil, err
		}
		if clause != "" {
			// Only add the operator if it's not the first group
			if i > 0 {
				clauses = append(clauses, string(group.Operator))
			}
			clauses = append(clauses, clause)
			args = append(args, groupArgs...)
			currentOffset += len(groupArgs)
		}
	}

	if len(clauses) == 0 {
		return "", nil, nil
	}

	return fmt.Sprintf("WHERE %s", strings.Join(clauses, " ")), args, nil
}
