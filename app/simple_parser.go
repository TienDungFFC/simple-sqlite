// this parser is used to parse simple query string in this challenge
package main

import (
	"fmt"
	"strings"
)

type Select struct {
	SelectExpr []string
	From       string
	Where      []string
}

func (s Select) isStatement() bool {
	return true
}
func SelectStatementParse(query string) (Statement, error) {
	selectStmt := Select{}
	qParts := strings.Split(strings.ToLower(query), " ")
	if qParts[0] != "select" {
		return nil, fmt.Errorf("This statement is not a select statement")
	}

	exprs := make([]string, 0)
	fromIdx := 0

	for idx, ex := range qParts[1:] {
		if ex == "from" {
			fromIdx = idx
			break
		}
		exprs = append(exprs, ex)
	}
	exprs = strings.Split(strings.Join(exprs, " "), ",")
	for i, expr := range exprs {
		exprs[i] = strings.TrimSpace(expr)
	}

	selectStmt.SelectExpr = exprs
	fromIdx += 2
	selectStmt.From = qParts[fromIdx]

	whereIdx := fromIdx + 2
	if whereIdx < len(qParts) {
		where := make([]string, 0)
		where = append(where, qParts[whereIdx:]...)
		selectStmt.Where = where
	}

	// only one table
	return selectStmt, nil
}
