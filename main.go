package main

import (
	"log"
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
		if conf.DiffType == "schema" {
			isUpdate = SchemaDiff(srcDB, dstDB, conf)
		} else {
			isUpdate = DataDiff(srcDB, dstDB, conf)
		}
		if isUpdate {
			log.Printf("Generated Update Script..%s!\n", conf.Output)
		} else {
			log.Println("Not Found Diff! ..BYE~!")
		}
	} else {
		MakeDoc(srcDB, conf)
	}
}
