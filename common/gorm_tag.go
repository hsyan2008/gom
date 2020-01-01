package common

import (
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

	if len(col.Default) >= 4 && strings.HasPrefix(col.Default, "''") && strings.HasSuffix(col.Default, "''") {
		col.Default = col.Default[1 : len(col.Default)-1]
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

	nstr := "type:" + strings.ToLower(DB().SQLType(col))
	res = append(res, nstr)

	if len(res) > 0 {
		return "gorm:\"" + strings.Join(res, ";") + "\""
	}

	return ""
}
