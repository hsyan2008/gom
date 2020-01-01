package common

import (
	"strings"

	"xorm.io/core"
)

type model struct {
	StructName string            //生成的struct名字
	TableName  string            //数据库里的表名，也是文件名
	Imports    map[string]string //需要导入的package列表
	Fields     []modelField      //对应的目标代码格式
}

func NewModel(table *core.Table) (m model) {
	m = model{
		StructName: core.LintGonicMapper.Table2Obj(table.Name),
		TableName:  table.Name,
		Imports:    make(map[string]string),
	}
	//postgres里可能存在重复字段
	var fieldMap = map[string]bool{}
	for _, column := range table.Columns() {
		if fieldMap[column.Name] {
			continue
		}
		fieldMap[column.Name] = true
		f := NewModelField(table, column)
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

func NewModelField(table *core.Table, column *core.Column) (f modelField) {
	f = modelField{
		FieldName:  core.LintGonicMapper.Table2Obj(column.Name),
		ColumnName: column.Name,
	}

	f.Type, f.Imports = getTypeAndImports(table, column)

	if column.Comment != "" {
		f.Comment = "// " + column.Comment
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
