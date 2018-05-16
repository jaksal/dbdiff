package main

import "fmt"
import "sort"

// Row : db row data..
type Row struct {
	Cols map[string]int
	Data []interface{}
}

// Get : get row form column name
func (r *Row) Get(key string) interface{} {
	if idx, exist := r.Cols[key]; exist {
		if r.Data[idx] != nil {
			return r.Data[idx]
		}
		return ""
	}
	return nil
}

// GetKeyValue : return keyvalue store
func (r *Row) GetKeyValue(keys []string) []*KeyValue {
	var result []*KeyValue

	if keys != nil {
		for _, k := range keys {
			result = append(result, &KeyValue{
				Key:   k,
				Value: r.Data[r.Cols[k]],
			})
		}
	} else {
		for i := 0; i < len(r.Data); i++ {
			result = append(result, &KeyValue{
				Key:   r.getKeyStr(i),
				Value: r.Data[i],
			})
		}
	}
	return result
}

func (r *Row) getKeyStr(idx int) string {
	for k, v := range r.Cols {
		if v == idx {
			return k
		}
	}
	return ""
}

// GetPrimaryKey : get primary key data.
func (r *Row) GetPrimaryKey(keys []string) string {
	var result string
	for _, k := range keys {
		result += fmt.Sprintf("%v|", r.Get(k))
	}
	return result
}

func (r *Row) compare(dst *Row) []string {
	var changeCols []string
	for i := 0; i < len(r.Data); i++ {
		if r.Data[i] != dst.Data[i] {
			changeCols = append(changeCols, r.getKeyStr(i))
		}
	}
	return changeCols
}

func reMapRows(src []*Row, keys []string) map[string]*Row {
	srcLists := make(map[string]*Row)
	for _, s := range src {
		key := s.GetPrimaryKey(keys)
		srcLists[key] = s
	}
	return srcLists
}

func checkIgnoreColumn(src, dst []string) bool {
	cnt := 0
	for _, ic := range src {
		for _, cc := range dst {
			if ic == cc {
				cnt++
				break
			}
		}
	}
	return cnt != len(dst)
}

// CompareRows : compare rows
func CompareRows(src, dst []*Row, keys []string, ignoreColumns []string) []*DiffRow {
	srcMap := reMapRows(src, keys)
	dstMap := reMapRows(dst, keys)

	var result []*DiffRow
	for k, s := range srcMap {
		// get dst..
		d, exist := dstMap[k]
		if !exist {
			continue
		}
		// compare.. rows..
		if changeCols := s.compare(d); changeCols != nil {
			if checkIgnoreColumn(ignoreColumns, changeCols) {
				result = append(result, &DiffRow{
					Compare: UPDATE,
					Keys:    s.GetKeyValue(keys),
					Data:    s.GetKeyValue(changeCols),
				})
			}
		}

		// delete..
		delete(srcMap, k)
		delete(dstMap, k)
	}

	// src insert..
	for _, s := range srcMap {
		result = append(result, &DiffRow{
			Compare: INSERT,
			Keys:    s.GetKeyValue(keys), // for sort
			Data:    s.GetKeyValue(nil),
		})
	}

	// dst drop..
	for _, d := range dstMap {
		result = append(result, &DiffRow{
			Compare: DELETE,
			Keys:    d.GetKeyValue(keys),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GetKey() < result[j].GetKey()
	})

	return result
}
