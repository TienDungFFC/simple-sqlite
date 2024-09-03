package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	interiorTable byte = 0x05
	leafTable     byte = 0x0d
)

var (
	TYPE_NULL  = 0
	TYPE_INT8  = 1
	TYPE_INT32 = 4
)

func (db *Database) HandleCommand(command string) {
	switch command {
	case ".dbinfo":
		db.GetInfo(os.Stdout)
	case ".tables":
		db.GetTables(os.Stdout)
	default:
		// default is select statement
		stmt, err := SelectStatementParse(command)
		if err != nil {
			panic(err)
		}
		db.Statement = stmt
		db.HandleStatement()
	}
}

func (db *Database) GetInfo(w io.Writer) {
	fmt.Printf("database page size: %v", db.Info.PageSize)
	fmt.Printf("number of tables: %v", db.Info.NumberOfCells)
	fmt.Fprintf(w, "database page size: %v", db.Info.PageSize)
	fmt.Fprintf(w, "number of tables: %v", db.Info.NumberOfCells)

}

func (db *Database) GetTables(w io.Writer) {
	var t string
	for _, table := range db.SchemaTables {
		if table.name != "sqlite_sequence" {
			t += table.name + " "
		}
	}
	fmt.Println(t[:len(t)-1])
	fmt.Fprint(w, t)
}

func (db *Database) ReadDb() {
	header := make([]byte, 100)

	_, err := db.f.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	var pageSize uint16
	if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
		fmt.Println("Failed to read integer:", err)
		return
	}
	db.Info.PageSize = pageSize

	var pageCount uint32
	if err := binary.Read(bytes.NewReader(header[28:32]), binary.BigEndian, &pageCount); err != nil {
		fmt.Println("Failed to read integer:", err)
		return
	}
	db.Info.PageCount = pageCount
	pageContent := make([]byte, pageSize)
	_, err = db.f.Read(pageContent[100:])
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(pageContent)
	reader.Seek(100, io.SeekStart)

	PageHeader, _ := readHeader(reader)

	switch PageHeader.pageType {
	case interiorTable:
		fmt.Println("interior table")
	case leafTable:
		cellPointerArray := make([]uint16, 0)
		var cellPointer uint16
		for i := 0; i < int(PageHeader.numberOfCells); i++ {
			binary.Read(reader, binary.BigEndian, &cellPointer)
			cellPointerArray = append(cellPointerArray, cellPointer)
		}
		schemaTable := make([]*SchemaTable, 0)
		for _, cellPointer := range cellPointerArray {
			reader.Seek(int64(cellPointer), io.SeekStart)
			_, _ = readVarint(reader)
			_, _ = readVarint(reader)
			totalHeaderSize, offset := readVarint(reader)
			colTypes := make([]uint64, 0)
			for offset < int(totalHeaderSize) {
				columnType, m := readVarint(reader)
				colTypes = append(colTypes, columnType)
				offset += m
			}
			colData := make([]any, 0)
			for _, col := range colTypes {

				switch col {
				case uint64(TYPE_INT8):
					val, _ := reader.ReadByte()
					colData = append(colData, int(val))
				default:
					if col >= 13 && col%2 != 0 {
						size := (col - 13) / 2
						content := make([]byte, size)
						reader.Read(content)
						colData = append(colData, content)
					}
				}
			}
			rootPage, _ := colData[3].(int)
			schema := SchemaTable{
				typ:      string(colData[0].([]byte)),
				name:     string(colData[1].([]byte)),
				tblName:  string(colData[2].([]byte)),
				rootPage: rootPage,
				sql:      string(colData[4].([]byte)),
			}
			schemaTable = append(schemaTable, &schema)
		}
		db.SchemaTables = schemaTable
	}
	db.Info.NumberOfCells = PageHeader.numberOfCells
}

func readHeader(reader *bytes.Reader) (*PageHeader, int) {
	header := PageHeader{}
	header.pageType, _ = reader.ReadByte()
	binary.Read(reader, binary.BigEndian, &header.startOfFirstFree)
	binary.Read(reader, binary.BigEndian, &header.numberOfCells)
	binary.Read(reader, binary.BigEndian, &header.startOfCellContent)
	binary.Read(reader, binary.BigEndian, &header.numberOfFragmentedFreeBytes)
	byteRead := 8
	return &header, byteRead
}

type PageHeader struct {
	pageType                    uint8
	startOfFirstFree            uint16
	numberOfCells               uint16
	startOfCellContent          uint16
	numberOfFragmentedFreeBytes uint8
	rightMostPointer            uint32
}

func (db *Database) HandleStatement() {
	switch stmt := db.Statement.(type) {
	case *Select:
		db.HandleSelectStatement(stmt)
	}
}

func (db *Database) ReadPayload(reader *bytes.Reader, colTypes []uint64) []any {
	colData := make([]any, 0)
	for _, col := range colTypes {
		switch col {
		case uint64(TYPE_NULL):
			colData = append(colData, nil)
		case uint64(TYPE_INT8):
			val, _ := reader.ReadByte()
			colData = append(colData, int(val))
		case uint64(TYPE_INT32):
			var val32 int
			binary.Read(reader, binary.BigEndian, &val32)
			colData = append(colData, val32)
		default:
			if col >= 13 && col%2 != 0 {
				size := (col - 13) / 2
				content := make([]byte, size)
				reader.Read(content)
				colData = append(colData, content)
			}
		}
	}
	return colData
}

func (db *Database) HandleSelectStatement(stmt *Select) {
	table := FindTable(db.SchemaTables, stmt.From)
	colTables := CreateStatementParse(table.sql)
	offset := (table.rootPage - 1) * int(db.Info.PageSize)
	db.f.Seek(int64(offset), io.SeekStart)
	pageContent := make([]byte, db.Info.PageSize)
	_, err := db.f.Read(pageContent)
	if err != nil {
		fmt.Println("Read error: ", err)
	}

	var countTable bool
	for _, expr := range stmt.SelectExpr {
		if strings.Contains("count(*)", expr) {
			countTable = true
			break
		}
	}
	reader := bytes.NewReader(pageContent)
	PageHeader, _ := readHeader(reader)
	if countTable {
		fmt.Println(PageHeader.numberOfCells)
	} else {
		cellPointerArray := make([]uint16, 0)
		var cellPointer uint16
		for i := 0; i < int(PageHeader.numberOfCells); i++ {
			binary.Read(reader, binary.BigEndian, &cellPointer)
			cellPointerArray = append(cellPointerArray, cellPointer)
		}
		for i := 0; i < int(PageHeader.numberOfCells); i++ {
			binary.Read(reader, binary.BigEndian, &cellPointer)
			cellPointerArray = append(cellPointerArray, cellPointer)
		}
		for _, cellPointer := range cellPointerArray {
			reader.Seek(int64(cellPointer), io.SeekStart)
			_, _ = readVarint(reader)
			_, _ = readVarint(reader)
			totalHeaderSize, offset := readVarint(reader)
			colTypes := make([]uint64, 0)
			for offset < int(totalHeaderSize) {
				columnType, m := readVarint(reader)
				colTypes = append(colTypes, columnType)
				offset += m
			}
			data := db.ReadPayload(reader, colTypes)
			recordData := make(map[string]string, 0)
			if len(data) > 0 {

				for i, col := range colTables {
					recordData[col] = fmt.Sprintf("%s", data[i])
				}

				for i, expr := range stmt.SelectExpr {
					if record, ok := recordData[expr]; ok {
						if i == len(stmt.SelectExpr)-1 {
							fmt.Println(record)
						} else {
							fmt.Print(record)
						}
					}
					if i < len(stmt.SelectExpr)-1 {
						fmt.Print("|")
					}
				}

			}

		}
	}

}
