package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Column : mysql Column desc.
type Column struct {
	Idx     int
	Name    string
	Type    string
	Null    string
	Default string
	Extra   string
	Comment string

	colidx int
	flag   int // 1=drop, 2=add , 3=update , 4=move
}

func (c *Column) String() string {
	return fmt.Sprintf("%d. name=%s type=%s null=%s default=%s extra=%s comment=%s",
		c.Idx, c.Name, c.Type, c.Null, c.Default, c.Extra, c.Comment)
}

func (c *Column) compare(dst *Column) bool {
	if c.Type != dst.Type ||
		c.Null != dst.Null ||
		c.Default != dst.Default ||
		c.Extra != dst.Extra ||
		c.Comment != dst.Comment {
		return false
	}
	return true
}

// GetSQL : generate alter table column info
func (c *Column) GetSQL() string {
	var result string

	result =  "`" + c.Name + "`" + " " + c.Type
	if c.Null == "YES" {
		result += " NULL"
	} else {
		result += " NOT NULL"
	}
	if c.Default != "" {
		result += " DEFAULT " + c.Default
	}
	if c.Extra != "" {
		result += " " + c.Extra
	}
	if c.Comment != "" {
		result += " COMMENT '" + CommentReplace(c.Comment) + "'"
	}
	return result
}

// Index db index desc.
type Index struct {
	Name     string
	Cols     []string
	IsUnique bool

	flag bool
}

func (i *Index) String() string {
	return fmt.Sprintf("name=%s cols=%+v is_unique=%t", i.Name, i.Cols, i.IsUnique)
}

func (i *Index) compare(dst *Index) bool {
	if len(i.Cols) != len(dst.Cols) {
		return false
	}
	for t := 0; t < len(i.Cols); t++ {
		if i.Cols[t] != dst.Cols[t] {
			return false
		}
	}
	if i.IsUnique != dst.IsUnique {
		return false
	}
	return true
}

func (i *Index) getAddSQL() string {
	var result string

	if i.Name == "PRIMARY" {
		result = "ADD PRIMARY KEY ("
	} else {
		if i.IsUnique {
			result = "ADD UNIQUE INDEX " + i.Name + " ("
		} else {
			result = "ADD INDEX " + i.Name + " ("
		}
	}
	for i, c := range i.Cols {
		if i > 0 {
			result += ","
		}
		result += c
	}

	return result + ")"
}
func (i *Index) getDropSQL() string {
	var result string

	if i.Name == "PRIMARY" {
		result = "DROP PRIMARY KEY"
	} else {
		result = "DROP INDEX " + i.Name
	}
	return result
}

// Table Table desc..
type Table struct {
	Name      string
	Cols      map[string]*Column
	Indexs    map[string]*Index
	Comment   string
	Engine    string
	Collation string
	Script    string
}

func (t *Table) String() string {
	return fmt.Sprintf("name=%s cols=%+v indexs=%+v comment=%s engine=%s collation=%s",
		t.Name, t.Cols, t.Indexs, t.Comment, t.Engine, t.Collation)
}

// GetPrimaryKey : get primary key column list
func (t *Table) GetPrimaryKey() []string {
	return t.Indexs["PRIMARY"].Cols
}

// GetColumn get column name from idx
func (t *Table) GetColumn(idx int) *Column {
	for _, c := range t.Cols {
		if c.Idx == idx {
			return c
		}
	}
	log.Fatalln("not found column idx", idx, t.Cols)
	return nil
}

// ParseTableDesc parse table object from create table script
/* CREATE TABLE `test` (
  `tsn` int(11) NOT NULL AUTO_INCREMENT,
  `tname` varchar(50) NOT NULL,
  `tvalue` int(11) unsigned zerofill DEFAULT '00000000001' COMMENT '123',
  `tdate` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '12321312',
  PRIMARY KEY (`tsn`),
  UNIQUE KEY `tname` (`tname`),
  KEY `tvalue` (`tvalue`)
) COMMENT='test 123' <nil>*/
func ParseTableDesc(desc string) *Table {
	return nil
}

// GetTableInfo get table desc .
func (d *DB) GetTableInfo(name string) *Table {
	// get table info..
	result := &Table{Name: name}
	{
		query := fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", name)
		data, err := d.GetData(query, nil)
		if err != nil {
			log.Fatalln("get data error", query, err)
			return nil
		}
		if len(data) == 0 {
			log.Println("emtpy table info", name)
			return nil
		}

		result.Comment = data[0].Get("Comment").(string)
		result.Engine = data[0].Get("Engine").(string)
		result.Collation = data[0].Get("Collation").(string)
	}

	// get column info..
	{
		query := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", name)
		data, err := d.GetData(query, nil)
		if err != nil {
			log.Fatalln("get data error", query, err)
			return nil
		}
		if len(data) == 0 {
			log.Println("emtpy column info", name)
			return nil
		}
		result.Cols = make(map[string]*Column)

		idx := 0
		for _, c := range data {
			col := &Column{
				Idx:     idx,
				Name:    c.Get("Field").(string),
				Type:    c.Get("Type").(string),
				Null:    c.Get("Null").(string),
				Default: c.Get("Default").(string),
				Extra:   c.Get("Extra").(string),
				Comment: c.Get("Comment").(string),
			}
			
			// 컬럼 타입이 숫자일 경우 int(10) -> int 로 변경.
			// DB 버전에 따라 다르게 보이는 이슈.
			// zerofill 옵션이 아닐경우에는 int 로 통일.
			if strings.HasPrefix(col.Type, "tinyint") ||
				strings.HasPrefix(col.Type, "smallint") ||
				strings.HasPrefix(col.Type, "mediumint") ||
				strings.HasPrefix(col.Type, "int") ||
				strings.HasPrefix(col.Type, "bigint") ||
				strings.HasPrefix(col.Type, "float") ||
				strings.HasPrefix(col.Type, "double") ||
				strings.HasPrefix(col.Type, "decimal") {
				col.Type = strings.Split(col.Type, "(")[0]
			}
			
			result.Cols[col.Name] = col
			idx++
		}
	}

	// get index..
	{
		query := fmt.Sprintf("SHOW INDEX FROM %s", name)
		data, err := d.GetData(query, nil)
		if err != nil {
			log.Fatalln("get data error", query, err)
			return nil
		}
		result.Indexs = make(map[string]*Index)
		for _, i := range data {
			key := i.Get("Key_name").(string)
			idx, exist := result.Indexs[key]
			if !exist {
				idx = &Index{Name: key}
				if i.Get("Non_unique").(string) == "0" {
					idx.IsUnique = true
				}
				result.Indexs[key] = idx
			}
			idx.Cols = append(idx.Cols, i.Get("Column_name").(string))
		}
	}

	// get create script
	{
		var err error
		result.Script, err = d.GetScript(TABLE, name)
		if err != nil {
			log.Fatalln("get table data error", name, err)
			return nil
		}
	}
	return result
}

func (t *Table) getBeforeColumn(idx int) string {
	if idx == 0 {
		return " FIRST"
	}
	return " AFTER " + t.GetColumn(idx-1).Name
}

func (t *Table) compare(d *Table) string {

	var updates []string

	var srcTable, dstTable []string

	// check drop cols
	colidx := 0
	for i := 0; i < len(d.Cols); i++ {
		dc := d.GetColumn(i)
		if _, exist := t.Cols[dc.Name]; !exist {
			sql := "DROP COLUMN " + dc.Name
			updates = append(updates, sql)

			dc.flag = 1
		} else {
			dc.colidx = colidx
			dstTable = append(dstTable, dc.Name)

			colidx++
		}
	}

	// check add and change cols.
	colidx = 0
	for i := 0; i < len(t.Cols); i++ {
		sc := t.GetColumn(i)
		if dc, exist := d.Cols[sc.Name]; !exist {
			sc.flag = 2
		} else {
			if !sc.compare(dc) {
				sc.flag = 3
			}
			sc.colidx = colidx
			srcTable = append(srcTable, sc.Name)

			colidx++
		}
	}

	// log.Println("src table", srcTable)
	// log.Println("dst table", dstTable)

	for i := 0; i < len(srcTable); {
		if srcTable[i] != dstTable[i] {
			// log.Println("diff column", i, srcTable[i], dstTable[i])

			sidx := t.Cols[dstTable[i]].colidx
			didx := d.Cols[srcTable[i]].colidx

			var start, end int
			if t.Cols[srcTable[i]].flag == 4 {
				start = i
				end = sidx
			} else {
				if sidx > didx {
					start = i
					end = sidx
				} else {
					start = didx
					end = i
				}
			}

			ch := dstTable[start]
			t.Cols[ch].flag = 4

			// log.Println("target idx", sidx, didx, start, end, ch)

			if start > end {
				dstTable = append(dstTable[:end],
					append([]string{ch},
						append(dstTable[end:start], dstTable[start+1:]...)...)...)
			} else {
				dstTable = append(dstTable[:start],
					append(dstTable[start+1:end+1],
						append([]string{ch}, dstTable[end+1:]...)...)...)
			}

			// log.Println("dst table ==>", dstTable)

			// log.Println("dst table #2", t2)
			if srcTable[i] == dstTable[i] {
				i++
			}

		} else {
			// log.Println("same col", i, srcTable[i])
			i++
		}
	}

	// getter add and update cols.
	for i := 0; i < len(t.Cols); i++ {
		sc := t.GetColumn(i)
		if sc.flag == 2 {
			// add
			sql := "ADD COLUMN " + sc.GetSQL() + t.getBeforeColumn(sc.Idx)
			updates = append(updates, sql)
		} else if sc.flag == 3 || sc.flag == 4 {
			// update
			sql := fmt.Sprintf("CHANGE COLUMN %s %s", sc.Name, sc.GetSQL()) + t.getBeforeColumn(sc.Idx)
			updates = append(updates, sql)
		}
	}

	// check index
	{
		for k, s := range t.Indexs {
			d := d.Indexs[k]
			if d == nil {
				continue
			}
			if !s.compare(d) {
				updates = append(updates, d.getDropSQL())
				updates = append(updates, s.getAddSQL())
			}
			s.flag = true
			d.flag = true
		}
		for _, s := range t.Indexs {
			if s.flag != true {
				updates = append(updates, s.getAddSQL())
			}
		}
		for _, s := range d.Indexs {
			if s.flag != true {
				updates = append(updates, s.getDropSQL())
			}
		}
	}

	// check ..
	if t.Comment != d.Comment {
		updates = append(updates, fmt.Sprintf("COMMENT='%s'", CommentReplace(t.Comment)))
	}

	if len(updates) == 0 {
		return ""
	}

	result := fmt.Sprintf("ALTER TABLE %s\n\t", t.Name)
	result += strings.Join(updates, ",\n\t")
	/*
		for i, u := range updates {
			if i != 0 {
				result += ", "
			}
			result += u
		}
	*/
	return result + ";"
}

// GenerateCreateSQL ...
func (t *Table) GenerateCreateSQL() string {
	return t.Script + ";"
}

// GenerateDropSQL ...
func (t *Table) GenerateDropSQL() string {
	return fmt.Sprintf("DROP TABLE %s;", t.Name)
}

// DiffTable ...
type DiffTable struct {
	Compare CompareType
	Name    string
	Result  string
}

// CompareTables compare table schema....
func CompareTables(srcs, dsts map[string]*Table) []*DiffTable {
	// check alter...

	var result []*DiffTable

	for k, t := range srcs {
		// get dst..
		d, exist := dsts[k]
		if !exist {
			continue
		}

		if alter := t.compare(d); alter != "" {
			result = append(result, &DiffTable{
				Compare: UPDATE,
				Name:    k,
				Result:  alter,
			})
		}
		delete(srcs, k)
		delete(dsts, k)
	}

	// src insert..
	for k, s := range srcs {
		result = append(result, &DiffTable{
			Compare: INSERT,
			Name:    k,
			Result:  s.GenerateCreateSQL(),
		})
	}

	// dst drop..
	for k, d := range dsts {
		result = append(result, &DiffTable{
			Compare: DELETE,
			Name:    k,
			Result:  d.GenerateDropSQL(),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
