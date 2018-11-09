package main

import (
	"log"
	"strings"
	"time"
)

func MakeDoc(srcDB *DB, conf *Config) {
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

	} else if conf.DiffType == "md" {
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
				out.Println("{code}sql")
				script, err := srcDB.GetScript(TABLE, srcTable.Name)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
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
				out.Println("{code}sql")
				script, err := srcDB.GetScript(VIEW, v)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
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
				out.Println("{code}sql")
				script, err := srcDB.GetScript(FUNCTION, f)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
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
				out.Println("{code}sql")
				script, err := srcDB.GetScript(PROCEDURE, p)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
				out.Println("[goto procedure list...](#procedures)")
				out.Println("<div style='page-break-after: always;'></div>")
				out.Printf("\n\n")
			}
		}
	} else if conf.DiffType == "wiki" {
		// output confluence wiki file.
		out.Println("h1. Schema : ", srcDB.DBName)

		// set table list table.
		out.Println("h3. Tables\n")
		out.Println("||name||comments||")

		var srcTables []*Table
		for _, t := range srcTableNames {
			srcTable := srcDB.GetTableInfo(t)
			if srcTable == nil {
				log.Fatal("load table info error", err)
			}
			out.Printf("|[%s|#%s]|%s |\n", srcTable.Name, WikiReplace(srcTable.Name), srcTable.Comment)

			srcTables = append(srcTables, srcTable)
		}
		out.Printf("\n\n----\n")

		// set view list table.
		out.Println("h3. Views\n")
		out.Println("||name||")

		for _, v := range srcViewNames {
			out.Printf("|[%s|#%s]|\n", v, WikiReplace(v))
		}
		out.Printf("\n\n----\n")

		// set function list.
		// set table list table.
		out.Println("h3. Functions\n")
		out.Println("||name||comments||")

		for _, f := range srcFunctionNames {
			out.Printf("|[%s|#%s]|%s |\n", f, WikiReplace(f), srcDB.GetObjectComments(FUNCTION, f))
		}
		out.Printf("\n\n----\n")

		// set procedure list.
		out.Println("h3. Procedures\n")
		out.Println("||name||comments||")

		for _, p := range srcProcedureNames {
			out.Printf("|[%s|#%s]|%s |\n", p, WikiReplace(p), srcDB.GetObjectComments(PROCEDURE, p))
		}
		out.Printf("\n\n----\n")
		out.Printf("\n\n----\n")

		// generate table doc.
		{
			// detail table ..
			for _, srcTable := range srcTables {
				out.Printf("h3. %s\n", srcTable.Name)
				out.Printf("bq. Comment : %s\\\\", srcTable.Comment)
				out.Printf("Engine : %s\\\\", srcTable.Engine)
				out.Printf("Collation : %s\\\\", srcTable.Collation)

				// print column..
				out.Println("\nh5. Columns\n")
				out.Println("||name||type||null||default||extra||comment||")
				for i := 0; i < len(srcTable.Cols); i++ {
					col := srcTable.GetColumn(i)
					out.Printf("|%s|%s|%s |%s |%s |%s |\n", col.Name, col.Type, col.Null, col.Default, col.Extra, col.Comment)
				}
				// print index.
				out.Println("\nh5. Indexs\n")
				out.Println("||name||columns||isnull||")
				for _, idx := range srcTable.Indexs {
					out.Printf("|%s|%s |%t |\n", idx.Name, strings.Join(idx.Cols, ","), idx.IsUnique)
				}
				out.Println("\nh5. Create Script\n")
				out.Println("{code}")
				script, err := srcDB.GetScript(TABLE, srcTable.Name)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
				out.Println("[goto table list...|#Tables]")
				out.Printf("\n\n")
			}
		}

		// generate view
		{
			for _, v := range srcViewNames {
				out.Printf("h3. %s\n\n\n", v)
				out.Println("\nh5. Create Script\n")
				out.Println("{code}")
				script, err := srcDB.GetScript(VIEW, v)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
				out.Println("[goto view list...|#Views]")
				out.Printf("\n\n")
			}
		}

		// generate function
		{
			for _, f := range srcFunctionNames {
				out.Printf("h3. %s\n", f)
				out.Printf("bq. Comment : %s  \n\n\n", srcDB.GetObjectComments(FUNCTION, f))
				out.Println("\nh5. Create Script\n")
				out.Println("{code}")
				script, err := srcDB.GetScript(FUNCTION, f)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
				out.Println("[goto function list...|#Functions]")
				out.Printf("\n\n")
			}
		}

		// generate procedure
		{
			for _, p := range srcProcedureNames {
				out.Printf("h3. %s\n", p)
				out.Printf("bq. Comment : %s  \n\n\n", srcDB.GetObjectComments(PROCEDURE, p))
				out.Println("\nh5. Create Script\n")
				out.Println("{code}")
				script, err := srcDB.GetScript(PROCEDURE, p)
				if err != nil {
					log.Fatal(err)
				}
				out.Println(script)
				out.Printf("{code}\n")
				out.Println("[goto procedure list...|#Procedures]")
				out.Printf("\n\n")
			}
		}
	} else {
		log.Fatal("invalid diff_type")
	}
}
