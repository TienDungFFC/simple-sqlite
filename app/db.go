package main

import "os"

type Statement interface {
	isStatement() bool
}
type Database struct {
	f            *os.File
	SchemaTables []*SchemaTable
	Info         *DatabaseInfo
	Statement    Statement
}

type DatabaseInfo struct {
	PageSize      uint16
	PageCount     uint32
	NumberOfCells uint16
}
type SchemaTable struct {
	typ      string
	name     string
	tblName  string
	rootPage int
	sql      string
}

type Payload struct {
	// totalHeaderSize
}

func NewDatabase(dbFile *os.File) *Database {
	return &Database{f: dbFile, SchemaTables: []*SchemaTable{}, Info: &DatabaseInfo{}}
}

func FindTable(tables []*SchemaTable, tblName string) *SchemaTable {
	for _, t := range tables {
		if t.tblName == tblName {
			return t
		}
	}
	return nil
}
