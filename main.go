package main

import (
	"fmt"

	"github.com/hsyan2008/gom/common"
)

func main() {
	var err error
	err = common.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%#v\n", common.Configs())

	genTool := common.NewGenTool()
	err = genTool.Gen()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("done")
}
