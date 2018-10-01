package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	// read config.
	conf, err := ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	// open source db
	srcDB := &DB{}
	if err := srcDB.Open(conf.Source); err != nil {
		log.Fatal(err, conf.Source)
	}
	defer srcDB.Close()

	if conf.DiffType == "schema" || conf.DiffType == "data" {
		// open target db
		dstDB := &DB{}
		if err := dstDB.Open(conf.Target); err != nil {
			log.Fatal(err, conf.Target)
		}
		defer dstDB.Close()

		var isUpdate bool

		out := &Output{}
		if err := out.Init(conf.Output); err != nil {
			log.Fatal(err, conf.Output)
		}
		defer func() {
			out.Close(isUpdate)
		}()

		out.Printf("-- Create Date : %s\n", time.Now())

		if conf.DiffType == "schema" {
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

		} else {
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
		}
		if isUpdate {
			log.Printf("Generated Update Script..%s!\n", conf.Output)
		} else {
			log.Println("Not Found Diff! ..BYE~!")
		}
	} else {
		// make md file.
		out := &Output{}
		if err := out.Init(conf.Output); err != nil {
			log.Fatal(err, conf.Output)
		}
		defer out.Close(true)

		out.Printf("-- Create Date : %s --\n", time.Now().Format("2006-01-02T15:04:05"))

		// get table list.
		srcTableNames, err := srcDB.GetObjectList(TABLE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}
		// get view list..
		srcViewNames, err := srcDB.GetObjectList(VIEW, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}
		srcFunctionNames, err := srcDB.GetObjectList(FUNCTION, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}
		srcProcedureNames, err := srcDB.GetObjectList(PROCEDURE, conf.Include, conf.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		if conf.DiffType == "sql" {
			// output sql file
			out.Printf("-- DB : %s\n\n", srcDB.DBName)

			// generate table.
			{
				// detail table ..
				for _, srcTable := range srcTableNames {
					out.Printf("-- %s\n", srcTable)
					script, err := srcDB.GetScript(TABLE, srcTable)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script + ";")
					out.Printf("\n\n")
				}
			}

			// generate view
			{
				for _, v := range srcViewNames {
					out.Printf("-- %s\n", v)
					script, err := srcDB.GetScript(VIEW, v)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script + ";")
					out.Printf("\n\n")
				}
			}

			// generate function
			{
				for _, f := range srcFunctionNames {
					out.Printf("-- %s\n", f)
					out.Println("DELIMITER //")
					script, err := srcDB.GetScript(FUNCTION, f)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script + "//")
					out.Println("DELIMITER ;")
					out.Printf("\n\n")
				}
			}

			// generate procedure
			{
				for _, p := range srcProcedureNames {
					out.Printf("-- %s\n", p)
					out.Println("DELIMITER //")
					script, err := srcDB.GetScript(PROCEDURE, p)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script + "//")
					out.Println("DELIMITER ;")
					out.Printf("\n\n")
				}
			}

		} else {
			// output md file.
			out.Println("# Schema : ", srcDB.DBName)

			// set table list table.
			out.Println("## Tables\n")
			out.Println("name | comments")
			out.Println(":--- | :---")

			var srcTables []*Table
			for _, t := range srcTableNames {
				srcTable := srcDB.GetTableInfo(t)
				if srcTable == nil {
					log.Fatal("load table info error", err)
				}
				out.Printf("[%s](#%s) | %s\n", srcTable.Name, MDReplace(srcTable.Name), srcTable.Comment)

				srcTables = append(srcTables, srcTable)
			}
			out.Println("<div style='page-break-after: always;'></div>")
			out.Printf("\n\n---\n")

			// set view list table.
			out.Println("## Views\n")
			out.Println("name |")
			out.Println(":--- |")

			for _, v := range srcViewNames {
				out.Printf("[%s](#%s) |\n", v, MDReplace(v))
			}
			out.Println("<div style='page-break-after: always;'></div>")
			out.Printf("\n\n---\n")

			// set function list.
			// set table list table.
			out.Println("## Functions\n")
			out.Println("name | comments")
			out.Println(":--- | :---")

			for _, f := range srcFunctionNames {
				out.Printf("[%s](#%s) | %s\n", f, MDReplace(f), srcDB.GetObjectComments(FUNCTION, f))
			}
			out.Println("<div style='page-break-after: always;'></div>")
			out.Printf("\n\n---\n")

			// set procedure list.
			out.Println("## Procedures\n")
			out.Println("name | comments")
			out.Println(":--- | :---")

			for _, p := range srcProcedureNames {
				out.Printf("[%s](#%s) | %s\n", p, MDReplace(p), srcDB.GetObjectComments(PROCEDURE, p))
			}
			out.Println("<div style='page-break-after: always;'></div>")
			out.Printf("\n\n---\n")

			// generate table.
			{
				// detail table ..
				for _, srcTable := range srcTables {
					out.Printf("## %s\n", srcTable.Name)
					out.Printf("> Comment : %s  \n", srcTable.Comment)
					out.Printf("> Engine : %s  \n", srcTable.Engine)
					out.Printf("> Collation : %s  \n\n", srcTable.Collation)

					// print column..
					out.Println("\nColumns\n")
					out.Println("name | type | null | default | extra | comment")
					out.Println(":--- | :--- | :--- | :--- | :--- | :---")
					for i := 0; i < len(srcTable.Cols); i++ {
						col := srcTable.GetColumn(i)
						out.Printf("%s | %s | %s | %s | %s | %s\n", col.Name, col.Type, col.Null, col.Default, col.Extra, col.Comment)
					}
					// print index.
					out.Println("\nIndexs\n")
					out.Println("name | columns | isnull")
					out.Println(":--- | :--- | :---")
					for _, idx := range srcTable.Indexs {
						out.Printf("%s | %s | %t \n", idx.Name, strings.Join(idx.Cols, ","), idx.IsUnique)
					}
					out.Println("\nCreate Script\n")
					out.Println("```sql")
					script, err := srcDB.GetScript(TABLE, srcTable.Name)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script)
					out.Printf("```\n")
					out.Println("[goto table list...](#tables)")
					out.Println("<div style='page-break-after: always;'></div>")
					out.Printf("\n\n")
				}
			}

			// generate view
			{
				for _, v := range srcViewNames {
					out.Printf("## %s\n\n\n", v)
					out.Println("\nCreate Script\n")
					out.Println("```sql")
					script, err := srcDB.GetScript(VIEW, v)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script)
					out.Printf("```\n")
					out.Println("[goto view list...](#views)")
					out.Println("<div style='page-break-after: always;'></div>")
					out.Printf("\n\n")
				}
			}

			// generate function
			{
				for _, f := range srcFunctionNames {
					out.Printf("## %s\n", f)
					out.Printf("> Comment : %s  \n\n\n", srcDB.GetObjectComments(FUNCTION, f))
					out.Println("\nCreate Script\n")
					out.Println("```sql")
					script, err := srcDB.GetScript(FUNCTION, f)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script)
					out.Printf("```\n")
					out.Println("[goto function list...](#functions)")
					out.Println("<div style='page-break-after: always;'></div>")
					out.Printf("\n\n")
				}
			}

			// generate procedure
			{
				for _, p := range srcProcedureNames {
					out.Printf("## %s\n", p)
					out.Printf("> Comment : %s  \n\n\n", srcDB.GetObjectComments(PROCEDURE, p))
					out.Println("\nCreate Script\n")
					out.Println("```sql")
					script, err := srcDB.GetScript(PROCEDURE, p)
					if err != nil {
						log.Fatal(err)
					}
					out.Println(script)
					out.Printf("```\n")
					out.Println("[goto procedure list...](#procedures)")
					out.Println("<div style='page-break-after: always;'></div>")
					out.Printf("\n\n")
				}
			}
		}
	}
}
