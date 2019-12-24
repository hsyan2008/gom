package common

import (
	"path"
	"strings"

	"xorm.io/core"
)

var created = []string{"created_at"}
var updated = []string{"updated_at"}
var deleted = []string{"deleted_at"}

func InStringSlice(f string, a []string) bool {
	for _, s := range a {
		if f == s {
			return true
		}
	}
	return false
}

func getTypeAndImports(column *core.Column) (t string, imports map[string]string) {
	t = sqlType2TypeString(column.SQLType)

	imports = map[string]string{}
	if len(Configs().ReplaceType) > 0 {
		for k, v := range Configs().ReplaceType {
			//原类型
			if t == path.Base(k) {
				//获取新包名
				pkg := strings.TrimSuffix(v, path.Ext(v))
				//获取新类型名
				t = path.Base(v)
				imports[pkg] = pkg
			}
		}
	} else {
		if t == "time.Time" {
			imports["time"] = "time"
		}
	}

	return
}
