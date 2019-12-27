package common

import (
	"fmt"

	"xorm.io/core"
)

func GetJsonTag(table *core.Table, column *core.Column) string {
	if Configs().JSONOmitempty {
		return fmt.Sprintf(`json:"%s,omitempty"`, column.Name)
	}
	if InStringSlice(column.Name, Configs().JSONIgnoreField) {
		return `json:"-"`
	}
	return fmt.Sprintf(`json:"%s"`, column.Name)
}
