package main

import (
	"testing"

	"github.com/xwb1989/sqlparser"
)

func TestFullScan(t *testing.T) {
	stmt, _ := sqlparser.Parse("SELECT id, name FROM superheroes WHERE eye_color = 'Pink Eyes'")
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		DBSelectCmd("../superheroes.db", stmt)	
	}
}