package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// DataDiff ...
func DataDiff(srcDB, dstDB *DB, conf *Config) bool {
	var isUpdate bool

	out := &Output{}
	if err := out.Init(conf.Output); err != nil {
		log.Fatal(err, conf.Output)
	}
	defer func() {
		out.Close(isUpdate)
	}()

	out.Printf("-- Create Date : %s\n", time.Now())

	// diff data
	// get table list.
	tableList, err := srcDB.GetObjectList(TABLE, conf.Include, conf.Exclude)
	if err != nil {
		log.Fatal(err)
	}

	// check target db.
	for _, t := range tableList {

		// log.Printf("compare table.....%s", t)
		// check target db.
		srcTable := srcDB.GetTableInfo(t)
		if srcTable == nil {
			log.Fatal("load table info error", err)
		}
		// check primary key.
		if _, ok := srcTable.Indexs["PRIMARY"]; !ok {
			log.Println("not exist primary key", t)
			continue
		}

		dstTable := dstDB.GetTableInfo(t)
		if srcTable == nil {
			log.Fatal("load table info error", err)
		}

		if dd := srcTable.compare(dstTable); dd != "" {
			log.Fatal("table desc is different!", t, dd)
		}

		query := fmt.Sprintf("SELECT * FROM %s", t)
		// load src db rows
		srcData, err := srcDB.GetData(query, srcTable.Cols)
		if err != nil {
			log.Fatal(err)
		}

		// load dst db rows
		dstData, err := dstDB.GetData(query, dstTable.Cols)
		if err != nil {
			log.Fatal(err)
		}

		// Compare data...
		result := CompareRows(srcData, dstData, srcTable.GetPrimaryKey(), strings.Split(conf.IgnoreColumn, ","))
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE UPDATE TABLE DATA %s\n", t)
			out.Printf("-- ----------------------------------------- --\n")
		}
		for _, r := range result {
			out.Println(r.GenerateSQL(t))
			isUpdate = true
		}
		log.Printf("compare table %s ... count=%d\n", t, len(result))
	}
	if isUpdate {
		out.Printf("\n\nCOMMIT;")
	}

	return isUpdate
}
