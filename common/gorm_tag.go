package common

import (
	"fmt"
	"sort"
	"strings"

	"xorm.io/core"
)

func GetGormTag(table *core.Table, col *core.Column) string {
	isNameId := col.Name == "id"
	isIdPk := isNameId && sqlType2TypeString(col.SQLType) == "int64"

	var res []string
	res = append(res, "column:"+col.Name)

	if !col.Nullable {
		if !isIdPk {
			res = append(res, "not null")
		}
	}
	if col.IsPrimaryKey {
		res = append(res, "primary_key")
	}
	if col.Default == "''''" {
		col.Default = "''"
	}
	if col.Default != "" {
		res = append(res, "default "+col.Default)
	}
	if col.IsAutoIncrement {
		res = append(res, "AUTO_INCREMENT")
	}

	names := make([]string, 0, len(col.Indexes))
	for name := range col.Indexes {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		index := table.Indexes[name]
		var uistr string
		if index.Type == core.UniqueType {
			if len(index.Cols) > 1 {
				uistr += "unique_index:" + index.Name
			} else {
				uistr = "unique"
			}
		} else if index.Type == core.IndexType {
			uistr = "index"
			if len(index.Cols) > 1 {
				uistr += ":" + index.Name
			}
		}
		res = append(res, uistr)
	}

	if col.SQLType.Name == "TIMESTAMPZ" {
		col.SQLType.Name = "TIMESTAMPTZ"
	}

	nstr := "type:" + strings.ToLower(col.SQLType.Name)
	if col.Length != 0 {
		if col.Length2 != 0 {
			nstr += fmt.Sprintf("(%v,%v)", col.Length, col.Length2)
		} else {
			nstr += fmt.Sprintf("(%v)", col.Length)
		}
	} else if len(col.EnumOptions) > 0 { //enum
		nstr += "("
		opts := ""

		enumOptions := make([]string, 0, len(col.EnumOptions))
		for enumOption := range col.EnumOptions {
			enumOptions = append(enumOptions, enumOption)
		}
		sort.Strings(enumOptions)

		for _, v := range enumOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	} else if len(col.SetOptions) > 0 { //enum
		nstr += "("
		opts := ""

		setOptions := make([]string, 0, len(col.SetOptions))
		for setOption := range col.SetOptions {
			setOptions = append(setOptions, setOption)
		}
		sort.Strings(setOptions)

		for _, v := range setOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	}
	res = append(res, nstr)

	if len(res) > 0 {
		return "gorm:\"" + strings.Join(res, ";") + "\""
	}

	return ""
}
