package main

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ObjectType : mysql object name
type ObjectType int

// DB Object Type
const (
	TABLE ObjectType = iota + 1
	VIEW
	FUNCTION
	PROCEDURE
	TRIGGER
)

/*
type ConnString string
func (cs *ConnString) Scan(state fmt.ScanState, verb rune) error {
	token, err := state.Token(true, unicode.IsLetter)
	if err != nil {
		return err
	}
	*cs = ConnString(token)
	return nil
}
*/

// DB mysql db ..
type DB struct {
	DB *sql.DB

	DBName string
}

// Open : open mysql.
func (d *DB) Open(constr string) error {
	db, err := sql.Open("mysql", constr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	d.DB = db
	d.DBName = constr[strings.IndexByte(constr, '/')+1:]

	log.Println("open mysql db...", constr, d.DBName)

	return nil
}

// Close : close mysql
func (d *DB) Close() {
	if d != nil && d.DB != nil {
		d.DB.Close()
	}
}

// GetData : get db rows
func (d *DB) GetData(query string, columns map[string]*Column) ([]*Row, error) {

	// Execute the query
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Set column index.
	cols := make(map[string]int)
	for i, c := range columnNames {
		cols[c] = i
	}

	values := make([]sql.RawBytes, len(columnNames))
	scanArgs := make([]interface{}, len(columnNames))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result []*Row

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		rowData := &Row{Cols: cols}

		for idx, col := range values {
			var value interface{}
			if col != nil {
				if columns != nil {
					if strings.HasPrefix(columns[columnNames[idx]].Type, "tinyint") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "smallint") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "mediumint") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "int") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "bigint") {
						value, _ = strconv.Atoi(string(col))
					} else if strings.HasPrefix(columns[columnNames[idx]].Type, "float") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "double") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "decimal") {
						value, _ = strconv.ParseFloat(string(col), 64)
					} else if strings.HasPrefix(columns[columnNames[idx]].Type, "char") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "varchar") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "json") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "text") {
						value = mysqlEscape(col)
					} else if strings.HasPrefix(columns[columnNames[idx]].Type, "date") ||
						strings.HasPrefix(columns[columnNames[idx]].Type, "datetime") {
						value = string(col)
					} else {
						log.Fatalln("invalid data type", columns[columnNames[idx]].Name, columns[columnNames[idx]].Type)
					}
				} else {
					value = string(col)
				}
			}
			rowData.Data = append(rowData.Data, value)
		}

		result = append(result, rowData)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

var reg1 = regexp.MustCompile("DEFINER=[^ ]* ")
var reg2 = regexp.MustCompile("ENGINE=[^ ]* ")
var reg3 = regexp.MustCompile("AUTO_INCREMENT=[^ ]* ")
var reg4 = regexp.MustCompile("DEFAULT CHARSET=[^ ]* ")
var reg5 = regexp.MustCompile("ALGORITHM=[^ ]* ")
var reg6 = regexp.MustCompile("ROW_FORMAT=[^ ]* ")

func ignoreDefiner(result string) string {
	result = reg1.ReplaceAllString(result, "")
	// result = reg2.ReplaceAllString(result, "")
	result = reg3.ReplaceAllString(result, "")
	result = reg4.ReplaceAllString(result, "")
	result = reg5.ReplaceAllString(result, "")
	result = reg6.ReplaceAllString(result, "")
	result = strings.Replace(result, "SQL SECURITY DEFINER ", "", -1)

	return result
}

// GetScript : get create script
func (d *DB) GetScript(objectType ObjectType, objectName string) (string, error) {
	var query string
	switch objectType {
	case TABLE:
		query = "SHOW CREATE TABLE " + objectName
	case VIEW:
		query = "SHOW CREATE VIEW " + objectName
	case FUNCTION:
		query = "SHOW CREATE FUNCTION " + objectName
	case PROCEDURE:
		query = "SHOW CREATE PROCEDURE " + objectName
	case TRIGGER:
		query = "SHOW CREATE TRIGGER " + objectName
	}

	//log.Println(query)

	data, err := d.GetData(query, nil)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("not found data")
	}

	// log.Println(data[0])

	var result string
	switch objectType {
	case TABLE:
		result = data[0].Get("Create Table").(string)
	case VIEW:
		result = data[0].Get("Create View").(string)
	case FUNCTION:
		result = data[0].Get("Create Function").(string)
	case PROCEDURE:
		result = data[0].Get("Create Procedure").(string)
	case TRIGGER:
		result = data[0].Get("Create Trigger").(string)
	}

	//log.Println(result)

	return strings.TrimSpace(ignoreDefiner(result)), nil
}

// GetObjectComments get object comment
func (d *DB) GetObjectComments(objectType ObjectType, objectName string) string {
	var query string
	switch objectType {
	case TABLE:
		query = "SHOW TABLE STATUS WHERE Name='" + objectName + "'"
	case VIEW, TRIGGER:
		return ""
	case FUNCTION:
		query = "SHOW FUNCTION STATUS WHERE Db = DATABASE() AND Name='" + objectName + "'"
	case PROCEDURE:
		query = "SHOW PROCEDURE STATUS WHERE Db = DATABASE() AND Name='" + objectName + "'"
	}

	data, err := d.GetData(query, nil)
	if err != nil || len(data) == 0 {
		return ""
	}
	return data[0].Get("Comment").(string)
}

// GetObjectList get db object list.
func (d *DB) GetObjectList(objectType ObjectType, include string, exclude string) ([]string, error) {
	var query string
	switch objectType {
	case TABLE:
		query = "SHOW FULL TABLES WHERE TABLE_TYPE NOT LIKE 'VIEW'"
	case VIEW:
		query = "SHOW FULL TABLES WHERE TABLE_TYPE LIKE 'VIEW'"
	case FUNCTION:
		query = "SHOW FUNCTION STATUS WHERE Db = DATABASE()"
	case PROCEDURE:
		query = "SHOW PROCEDURE STATUS WHERE Db = DATABASE()"
	case TRIGGER:
		query = "SHOW TRIGGERS"
	}

	data, err := d.GetData(query, nil)
	if err != nil {
		return nil, err
	}

	if include != "" {
		include = strings.ToLower(include)
	}
	if exclude != "" {
		exclude = strings.ToLower(exclude)
	}

	var result []string
	for _, d := range data {
		var name string
		switch objectType {
		case TABLE, VIEW:
			name = d.Data[0].(string)
		case FUNCTION, PROCEDURE:
			name = d.Get("Name").(string)
		case TRIGGER:
			name = d.Get("Trigger").(string)
		}

		bIgnore := false
		if include != "" && strings.Contains(strings.ToLower(name), include) == false {
			bIgnore = true
		}
		if exclude != "" && strings.Contains(strings.ToLower(name), exclude) == true {
			bIgnore = true
		}

		if bIgnore == false {
			result = append(result, name)
		}
	}
	return result, nil
}

// GetTableList get table info list.
func (d *DB) GetTableList(include string, exclude string) (map[string]*Table, error) {
	tableNames, err := d.GetObjectList(TABLE, include, include)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*Table)
	for _, t := range tableNames {
		table := d.GetTableInfo(t)
		if table == nil {
			return nil, errors.New("load table info error")
		}
		result[t] = table
	}
	return result, nil
}

// GetScriptList get script info list.
func (d *DB) GetScriptList(objectType ObjectType, include string, exclude string) (map[string]string, error) {
	scriptNames, err := d.GetObjectList(objectType, include, include)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, t := range scriptNames {
		script, err := d.GetScript(objectType, t)
		if err != nil {
			return nil, err
		}
		result[t] = script
	}
	return result, nil
}
