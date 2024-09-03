// this parser is used to parse simple query string in this challenge
package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Select struct {
	SelectExpr []string
	From       string
	Where      []string
}

func (s *Select) isStatement() bool {
	return true
}
func SelectStatementParse(query string) (Statement, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if !strings.HasPrefix(query, "select ") {
		return nil, fmt.Errorf("this statement is not a select statement")
	}

	qParts := strings.Fields(query)
	fromIdx := indexOf(qParts, "from")
	if fromIdx == -1 {
		return nil, fmt.Errorf("'from' keyword not found")
	}

	exprs := qParts[1:fromIdx]
	selectExpr := splitAndTrim(strings.Join(exprs, " "))

	fromIdx++
	if fromIdx >= len(qParts) {
		return nil, fmt.Errorf("table name missing after 'from'")
	}
	table := qParts[fromIdx]

	whereClause := []string{}
	if fromIdx+1 < len(qParts) && qParts[fromIdx+1] == "where" {
		whereClause = qParts[fromIdx+2:]
	}

	selectStmt := &Select{
		SelectExpr: selectExpr,
		From:       table,
		Where:      whereClause,
	}

	return selectStmt, nil
}

func CreateStatementParse(stmt string) []string {
	re := regexp.MustCompile(`(?s)CREATE TABLE [^\(\)]+(?:\s*\((.*?)\))`)
	match := re.FindStringSubmatch(stmt)

	if len(match) > 1 {
		fields := match[1]
		fieldNames := extractFieldNames(fields)
		return fieldNames
	} else {
		fmt.Println("No CREATE TABLE statement found.")
	}
	return nil
}

func extractFieldNames(fields string) []string {
	var fieldNames []string

	lines := strings.Split(fields, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if index := strings.IndexAny(line, " \t"); index != -1 {
			fieldName := line[:index]
			fieldNames = append(fieldNames, fieldName)
		}
	}

	return fieldNames
}

func indexOf(parts []string, keyword string) int {
	for i, part := range parts {
		if part == keyword {
			return i
		}
	}
	return -1
}

func splitAndTrim(expr string) []string {
	exprs := strings.Split(expr, ",")
	for i, ex := range exprs {
		exprs[i] = strings.TrimSpace(ex)
	}
	return exprs
}
