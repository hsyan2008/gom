package common

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hsyan2008/hfw/common"
	"xorm.io/core"
)

type GenTool struct {
	targetDir   string
	packageName string
	tables      []*core.Table
	models      map[string]model
}

func NewGenTool() *GenTool {
	dir := Configs().TargetDir
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(common.GetAppPath(), dir)
	}
	if !common.IsExist(dir) {
		os.MkdirAll(dir, 0755)
	}
	return &GenTool{
		targetDir:   dir,
		packageName: filepath.Base(dir),
		models:      make(map[string]model),
	}
}

func (genTool *GenTool) getDBMetas() (err error) {
	genTool.tables, err = DBMetas(Configs().Tables, Configs().ExcludeTables, Configs().TryComplete)
	if err != nil {
		return
	}

	return nil
}

func (genTool *GenTool) genModels() {
	for _, table := range genTool.tables {
		model := NewModel(table)
		genTool.models[model.TableName] = model
	}

	return
}

func (genTool *GenTool) genFile() (err error) {
	for _, model := range genTool.models {
		log.Println("start gen table:", model.TableName)
		str := fmt.Sprintln("package", genTool.packageName)
		if len(model.Imports) > 0 {
			str += fmt.Sprintln("import (")
			for _, i := range model.Imports {
				str += fmt.Sprintf(`"%s"`, i)
			}
			str += fmt.Sprintln(")")
		}
		str += fmt.Sprintln("type", model.StructName, "struct {")
		for _, v := range model.Fields {
			str += fmt.Sprintln(v.FieldName, v.Type, v.Tag)
		}
		str += fmt.Sprintln("}")
		str += fmt.Sprintln("func (", model.StructName, ") TableName() string {")
		str += fmt.Sprintln(fmt.Sprintf("return `%s`", model.TableName))
		str += fmt.Sprintln("}")

		b, err := format.Source([]byte(str))
		if err != nil {
			return err
		}
		file := filepath.Join(genTool.targetDir, fmt.Sprintf("%s.go", model.TableName))
		err = ioutil.WriteFile(file, b, 0644)
		if err != nil {
			return err
		}
		log.Println("gen to:", file)
	}
	return
}

func (genTool *GenTool) Gen() (err error) {
	if err = InitDb(); err != nil {
		return
	}

	if err = genTool.getDBMetas(); err != nil {
		return
	}

	genTool.genModels()

	if err = genTool.genFile(); err != nil {
		return
	}

	return nil
}
