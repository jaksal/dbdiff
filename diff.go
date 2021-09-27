package main

import (
	"fmt"
	"log"
	"strconv"
)

// CompareType contents compare type
type CompareType int

// CompareType enum..
const (
	UPDATE CompareType = iota + 1
	INSERT
	DELETE
)

// KeyValue : key value store
type KeyValue struct {
	Key   string
	Value interface{}
}

func (kv *KeyValue) String() string {
	switch kv.Value.(type) {
	case int:
		return fmt.Sprintf("%20d", kv.Value.(int))
	case float64:
		return fmt.Sprintf("%20f", kv.Value.(float64))
	default:
		return kv.Value.(string)
	}
	// log.Fatalln("not found key type!", kv)
}

// DiffRow : diff row info.
type DiffRow struct {
	Compare CompareType // 1=alter , 2 = insert, 3 = delete
	Data    []*KeyValue
	Keys    []*KeyValue
}

// GetKey : get key data string
func (d *DiffRow) GetKey() string {
	var result string
	for _, k := range d.Keys {
		result += fmt.Sprintf("%s|", k)
	}
	return result
}

func getDataStr(v interface{}) string {
	switch t := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", v.(string))
	case int:
		return strconv.Itoa(v.(int))
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case nil:
		return "NULL"
	default:
		log.Fatalln("check type..", t)
	}
	return ""
}

// GenerateSQL : generate sql string
func (d *DiffRow) GenerateSQL(name string) string {
	switch d.Compare {
	case INSERT:
		var keys, values string
		for _, d := range d.Data {
			if len(keys) > 0 {
				keys += ","
				values += ","
			}
			keys += d.Key
			values += getDataStr(d.Value)
		}
		return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", name, keys, values)
	case UPDATE:
		var keys, values string
		for _, k := range d.Keys {
			if len(keys) > 0 {
				keys += " AND "
			}
			keys += fmt.Sprintf("%s=%s", k.Key, getDataStr(k.Value))
		}
		for _, d := range d.Data {
			if len(values) > 0 {
				values += ", "
			}
			values += fmt.Sprintf("%s=%s", d.Key, getDataStr(d.Value))
		}
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", name, values, keys)
	case DELETE:
		var keys string
		for _, k := range d.Keys {
			if len(keys) > 0 {
				keys += " AND "
			}
			keys += fmt.Sprintf("%s=%s", k.Key, getDataStr(k.Value))
		}
		return fmt.Sprintf("DELETE FROM %s WHERE %s;", name, keys)
	}
	return ""
}
