package common

import (
	"strings"

	"github.com/hsyan2008/gom/common/big_camel"
	"xorm.io/core"
)

type model struct {
	StructName string            //生成的struct名字
	TableName  string            //数据库里的表名，也是文件名
	Imports    map[string]string //需要导入的package列表
	Fields     []modelField      //对应的目标代码格式
}

func NewModel(table *core.Table) (m model, err error) {
	m = model{
		StructName: big_camel.Marshal(table.Name),
		TableName:  table.Name,
		Imports:    make(map[string]string),
	}
	for _, column := range table.Columns() {
		// fmt.Printf("%s%#v\n", table.Name, column)
		f, err := NewModelField(table, column)
		if err != nil {
			return m, err
		}
		for k, v := range f.Imports {
			m.Imports[k] = v
		}
		m.Fields = append(m.Fields, f)
	}
	return
}

//每个字段的内容
type modelField struct {
	FieldName  string            //struct里的字段名
	ColumnName string            //表里的字段名，也是json名
	Type       string            //struct的字段类型
	Imports    map[string]string //需要导入的package列表
	Tag        string
	Comment    string //备注，放到注释里
}

func NewModelField(table *core.Table, column *core.Column) (f modelField, err error) {
	f = modelField{
		FieldName:  big_camel.Marshal(column.Name),
		ColumnName: column.Name,
		Type:       sqlType2TypeString(column.SQLType),
		Imports:    getGoImports(column),
		Comment:    column.Comment,
	}
	tags := []string{}
	for _, v := range Configs().TagType {
		switch v {
		case "json":
			tags = append(tags, GetJsonTag(table, column))
		case "xorm":
			tags = append(tags, GetXormTag(table, column))
		case "gorm":
			tags = append(tags, GetGormTag(table, column))
		}
	}
	if len(tags) > 0 {
		f.Tag = "`" + strings.Join(tags, " ") + "`"
	}

	return
}
