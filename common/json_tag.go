package common

import (
	"fmt"

	"xorm.io/core"
)

func GetJsonTag(table *core.Table, column *core.Column) string {
	return fmt.Sprintf(`json:"%s"`, column.Name)
}

func GetJsonTagWithOmitEmpty(table *core.Table, column *core.Column) string {
	return fmt.Sprintf(`json:"%s,omitempty"`, column.Name)
}
