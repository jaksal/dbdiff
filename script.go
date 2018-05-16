package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode"
)

// DiffScript diff script
type DiffScript struct {
	Compare CompareType
	Type    ObjectType
	Name    string
	Script  string
}

func (ds *DiffScript) generateDropSQL() string {
	switch ds.Type {
	case TABLE:
		return fmt.Sprintf("DROP TABLE %s", ds.Name)
	case VIEW:
		return fmt.Sprintf("DROP VIEW %s", ds.Name)
	case FUNCTION:
		return fmt.Sprintf("DROP FUNCTION %s", ds.Name)
	case PROCEDURE:
		return fmt.Sprintf("DROP PROCEDURE %s", ds.Name)
	case TRIGGER:
		return fmt.Sprintf("DROP TRIGGER %s", ds.Name)
	}
	log.Fatal("invalid object type", ds.Type)
	return ""
}

// GenerateSQL generate update sql
func (ds *DiffScript) GenerateSQL() string {
	switch ds.Compare {
	case INSERT:
		switch ds.Type {
		case FUNCTION, PROCEDURE, TRIGGER:
			return ds.Script + RoutineDelimeter
		default:
			return ds.Script + ";"
		}
	case UPDATE:
		switch ds.Type {
		case VIEW:
			return strings.Replace(ds.Script, "CREATE", "CREATE OR REPLACE", -1) + ";"
		case FUNCTION, PROCEDURE, TRIGGER:
			return fmt.Sprintf("%s%s\n%s%s", ds.generateDropSQL(), RoutineDelimeter, ds.Script, RoutineDelimeter)
		}
	case DELETE:
		switch ds.Type {
		case FUNCTION, PROCEDURE, TRIGGER:
			return ds.generateDropSQL() + RoutineDelimeter
		default:
			return ds.generateDropSQL() + ";"
		}
	}
	log.Fatal("invalid compare type or compare", ds.Compare, ds.Type)
	return ""
}

func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

func compareImpl(src, dst string) bool {
	if src == dst {
		return true
	}

	src = stringMinifier(src)
	dst = stringMinifier(dst)

	return strings.EqualFold(src, dst)
	/*
		srcLine := strings.Split(src, "\n")
		dstLine := strings.Split(dst, "\n")
		if len(srcLine) != len(dstLine) {
			return false
		}
		for i := 0; i < len(srcLine); i++ {
			if strings.ToLower(strings.TrimSpace(srcLine[i])) !=
				strings.ToLower(strings.TrimSpace(dstLine[i])) {
				return false
			}
		}
		return true
	*/
}

// CompareScript two script
func CompareScript(objectType ObjectType, srcs, dsts map[string]string) []*DiffScript {

	var result []*DiffScript

	for k, ss := range srcs {
		// get dst..
		ds, exist := dsts[k]
		if !exist {
			continue
		}
		if !compareImpl(ss, ds) {
			result = append(result, &DiffScript{
				Compare: UPDATE,
				Type:    objectType,
				Name:    k,
				Script:  ss,
			})
			//log.Printf("%s is different\n%s\n%s", k, ss, ds)
		}

		delete(srcs, k)
		delete(dsts, k)
	}

	// src insert..
	for k, ss := range srcs {
		result = append(result, &DiffScript{
			Compare: INSERT,
			Type:    objectType,
			Name:    k,
			Script:  ss,
		})
		log.Printf("%s is new\n", k)
	}

	// dst drop..
	for k, ds := range dsts {
		result = append(result, &DiffScript{
			Compare: DELETE,
			Type:    objectType,
			Name:    k,
			Script:  ds,
		})
		log.Printf("%s is deleted\n", k)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
