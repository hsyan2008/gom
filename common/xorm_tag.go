package common

import (
	"sort"
	"strings"

	"xorm.io/core"
)

func GetXormTag(table *core.Table, col *core.Column) string {
	isNameId := col.Name == "id"
	isIdPk := isNameId && sqlType2TypeString(col.SQLType) == "int64"

	var res []string
	if !col.Nullable {
		if !isIdPk {
			res = append(res, "not null")
		}
	}

	if col.IsPrimaryKey {
		res = append(res, "pk")
	}

	if len(col.Default) >= 4 && strings.HasPrefix(col.Default, "''") && strings.HasSuffix(col.Default, "''") {
		col.Default = col.Default[1 : len(col.Default)-1]
	}
	if col.Default != "" {
		res = append(res, "default "+col.Default)
	}

	if col.IsAutoIncrement {
		res = append(res, "autoincr")
	}

	if col.SQLType.IsTime() && InStringSlice(col.Name, created) {
		res = append(res, "created")
	}

	if col.SQLType.IsTime() && InStringSlice(col.Name, updated) {
		res = append(res, "updated")
	}

	if col.SQLType.IsTime() && InStringSlice(col.Name, deleted) {
		res = append(res, "deleted")
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
			uistr = "unique"
		} else if index.Type == core.IndexType {
			uistr = "index"
		}
		if len(index.Cols) > 1 {
			uistr += "(" + index.Name + ")"
		}
		res = append(res, uistr)
	}

	res = append(res, DB().SQLType(col))

	if len(res) > 0 {
		return "xorm:\"" + strings.Join(res, " ") + "\""
	}

	return ""
}
