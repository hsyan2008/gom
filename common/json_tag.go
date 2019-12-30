package common

import (
	"fmt"

	"xorm.io/core"
)

func GetJsonTag(table *core.Table, column *core.Column) string {
	var omit = ""
	if Configs().JSONOmitempty {
		omit = ",omitempty"
	}

	var columnName = column.Name
	if InStringSlice(column.Name, Configs().JSONIgnoreField) ||
		InStringSlice(fmt.Sprintf("%s.%s", table.Name, column.Name), Configs().JSONIgnoreField) {
		columnName = "-"
	}

	return fmt.Sprintf(`json:"%s"%s`, columnName, omit)
}
