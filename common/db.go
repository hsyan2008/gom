package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/core"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitDb() (err error) {
	engine, err = xorm.NewEngine(Configs().Driver, Configs().Source)
	if err != nil {
		return
	}
	engine.SetLogLevel(core.LOG_WARNING)

	if err = engine.Ping(); err != nil {
		return
	}

	return nil
}

func DB() *xorm.Engine {
	return engine
}

// DBMetas
// t 指定表，不指定就所有表
func DBMetas(t ...string) (tables []*core.Table, err error) {
	//类似DBMetas，因为一次性获取，碰到pgsql自定义的类型，直接出错，通过下面的方法，可以过滤

	dialect := DB().Dialect()
	tmpTables, err := dialect.GetTables()
	if err != nil {
		return
	}
	for _, v := range tmpTables {
		if len(t) > 0 && !InStringSlice(v.Name, t) {
			continue
		}
		if err = loadTableInfo(v); err != nil {
			return
		}
		tables = append(tables, v)
	}

	return
}

func loadTableInfo(table *core.Table) error {
	colSeq, cols, err := DB().Dialect().GetColumns(table.Name)
	if err != nil {
		return err
	}
	for _, name := range colSeq {
		table.AddColumn(cols[name])
	}
	indexes, err := DB().Dialect().GetIndexes(table.Name)
	if err != nil {
		return err
	}
	table.Indexes = indexes

	for _, index := range indexes {
		for _, name := range index.Cols {
			if col := table.GetColumn(name); col != nil {
				col.Indexes[index.Name] = index.Type
			} else {
				return fmt.Errorf("Unknown col %s in index %v of table %v, columns %v", name, index.Name, table.Name, table.ColumnsSeq())
			}
		}
	}
	return nil
}

func sqlType2TypeString(st core.SQLType) string {
	t := core.SQLType2Type(st)
	s := t.String()
	if s == "[]uint8" {
		return "[]byte"

	}
	return s
}

func getGoImports(column *core.Column) map[string]string {
	imports := make(map[string]string)

	if sqlType2TypeString(column.SQLType) == "time.Time" {
		imports["time"] = "time"
	}

	return imports
}
