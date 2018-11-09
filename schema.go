package main

import (
	"log"
	"time"
)

// SchemaDiff ...
func SchemaDiff(srcDB, dstDB *DB, conf *Config) bool {
	var isUpdate bool

	out := &Output{}
	if err := out.Init(conf.Output); err != nil {
		log.Fatal(err, conf.Output)
	}
	defer func() {
		out.Close(isUpdate)
	}()

	out.Printf("-- Create Date : %s\n", time.Now())
	// diff table
	{
		log.Println("compare tables....")
		srcTableNames, err := srcDB.GetObjectList(TABLE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}
		srcTables := make(map[string]*Table)
		for _, t := range srcTableNames {
			table := srcDB.GetTableInfo(t)
			if table == nil {
				log.Fatal("load table info error", err)
			}
			srcTables[t] = table
		}

		dstTableNames, err := dstDB.GetObjectList(TABLE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}
		dstTables := make(map[string]*Table)
		for _, t := range dstTableNames {
			table := dstDB.GetTableInfo(t)
			if table == nil {
				log.Fatal("load table info error", table)
			}
			dstTables[t] = table
		}

		// Compare Table.
		result := CompareTables(srcTables, dstTables)
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE TABLE SCHEMA \n")
			out.Printf("-- ----------------------------------------- --\n")
		}
		for _, r := range result {
			out.Printf("\n-- %s ...\n", r.Name)
			out.Println(r.Result)
			isUpdate = true
		}
		if len(result) > 0 {
			isUpdate = true
		}
		log.Println("finish tables....", len(result))
	}
	// diff view..
	{
		log.Println("compare views....")
		srcViews, err := srcDB.GetScriptList(VIEW, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		dstViews, err := dstDB.GetScriptList(VIEW, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		// Compare View.
		result := CompareScript(VIEW, srcViews, dstViews)
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE VIEW SCHEMA \n")
			out.Printf("-- ----------------------------------------- --\n")
		}
		for _, r := range result {
			out.Printf("\n-- %s ...\n", r.Name)
			out.Println(r.GenerateSQL())
		}
		if len(result) > 0 {
			isUpdate = true
		}
		log.Println("finish views....", len(result))
	}

	// diff function
	{
		log.Println("compare function....")
		srcViews, err := srcDB.GetScriptList(FUNCTION, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		dstViews, err := dstDB.GetScriptList(FUNCTION, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		// Compare View.
		result := CompareScript(FUNCTION, srcViews, dstViews)
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE FUNCTION SCHEMA \n")
			out.Printf("-- ----------------------------------------- --\n")
			out.Printf("DELIMITER %s\n", RoutineDelimeter)
		}
		for _, r := range result {
			out.Printf("\n-- %s ...\n", r.Name)
			out.Println(r.GenerateSQL())
		}
		if len(result) > 0 {
			out.Printf("\nDELIMITER ;")
			isUpdate = true
		}
		log.Println("finish function....", len(result))
	}

	// diff procedure
	{
		log.Println("compare procedure....")
		srcViews, err := srcDB.GetScriptList(PROCEDURE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		dstViews, err := dstDB.GetScriptList(PROCEDURE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		// Compare View.
		result := CompareScript(PROCEDURE, srcViews, dstViews)
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE PROCEDURE SCHEMA \n")
			out.Printf("-- ----------------------------------------- --\n")
			out.Printf("DELIMITER %s\n", RoutineDelimeter)
		}
		for _, r := range result {
			out.Printf("\n-- %s ...\n", r.Name)
			out.Println(r.GenerateSQL())
		}
		if len(result) > 0 {
			out.Printf("\nDELIMITER ;")
			isUpdate = true
		}
		log.Println("finish procedure....", len(result))
	}

	// diff trigger
	{
		log.Println("compare trigger....")
		srcViews, err := srcDB.GetScriptList(TRIGGER, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		dstViews, err := dstDB.GetScriptList(TRIGGER, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		// Compare trigger.
		result := CompareScript(TRIGGER, srcViews, dstViews)
		if len(result) > 0 {
			out.Printf("\n-- ----------------------------------------- --\n")
			out.Printf("-- GENERATE TRIGGER SCHEMA \n")
			out.Printf("-- ----------------------------------------- --\n")
			out.Printf("DELIMITER %s\n", RoutineDelimeter)
		}
		for _, r := range result {
			out.Printf("\n-- %s ...\n", r.Name)
			out.Println(r.GenerateSQL())
		}
		if len(result) > 0 {
			out.Printf("\nDELIMITER ;")
			isUpdate = true
		}
		log.Println("finish trigger....", len(result))
	}
	return isUpdate
}
