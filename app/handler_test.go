package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintPageSize(t *testing.T) {
	databaseFile, err := os.Open("../sample.db")
	if err != nil {
		panic(err)
	}
	defer databaseFile.Close()
	var buf = new(bytes.Buffer)

	db := NewDatabase(databaseFile)
	db.ReadDb()
	expected := "database page size: 4096"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}

func TestNumberOfTables(t *testing.T) {
	databaseFile, err := os.Open("../sample.db")
	if err != nil {
		panic(err)
	}
	defer databaseFile.Close()
	var buf = new(bytes.Buffer)

	db := NewDatabase(databaseFile)
	db.ReadDb()
	db.GetInfo(buf)
	expected := "number of tables: 3"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}

func TestPrintTable(t *testing.T) {
	databaseFile, err := os.Open("../sample.db")
	if err != nil {
		panic(err)
	}
	defer databaseFile.Close()
	var buf = new(bytes.Buffer)

	db := NewDatabase(databaseFile)
	db.ReadDb()
	db.GetTables(buf)
	expected := "apples"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}

func TestSelectStatement(t *testing.T) {
	databaseFile, err := os.Open("../sample.db")
	if err != nil {
		panic(err)
	}
	defer databaseFile.Close()
	var buf = new(bytes.Buffer)

	db := NewDatabase(databaseFile)
	db.ReadDb()
	db.HandleCommand("SELECT COUNT(*) FROM apples")
	expected := "4"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}

func TestSelectStatement2(t *testing.T) {
	sql := "SELECT id, name FROM companies WHERE country = 'eritrea'"
	stmt, _ := SelectStatementParse(sql)
	t.Error(stmt)
}

func TestHandleSelectStatement(t *testing.T) {
	databaseFile, err := os.Open("../sample.db")
	if err != nil {
		panic(err)
	}
	defer databaseFile.Close()

	db := NewDatabase(databaseFile)
	db.ReadDb()
	stmt, _ := SelectStatementParse("SELECT COUNT(*) FROM apples")
	s, _ := stmt.(Select)
	db.HandleSelectStatement(s)
	t.Errorf("!@3")
}
