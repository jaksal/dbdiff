package main

import (
	"fmt"
	"log"
	"os"
)

//Output . output file...
type Output struct {
	f *os.File
}

//Init ...
func (o *Output) Init(name string) error {
	if name != "" {
		name = CheckFilePath(name)

		var err error
		o.f, err = os.Create(name)
		if err != nil {
			return err
		}
		log.Println("open result file ", name)
	}
	return nil
}

//Close ..
func (o *Output) Close(isUpdate bool) {
	if o.f != nil {
		o.f.Close()

		if isUpdate == false {
			os.Remove(o.f.Name())
		}
	}
}

// Printf ...
func (o *Output) Printf(format string, v ...interface{}) {
	if o.f != nil {
		fmt.Fprintf(o.f, format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

// Println ...
func (o *Output) Println(v ...interface{}) {
	if o.f != nil {
		fmt.Fprintln(o.f, v...)
	} else {
		fmt.Println(v...)
	}
}
