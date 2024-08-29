package main

import (
	"fmt"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]
	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		databaseFile.Close()
	}()

	db := NewDatabase(databaseFile)
	fmt.Println("db created", db)
	db.ReadDb()
	db.HandleCommand(command)
}
