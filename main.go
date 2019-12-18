package main

import (
	"log"

	"github.com/hsyan2008/gom/common"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	if err := common.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := common.NewGenTool().Gen(); err != nil {
		log.Fatal(err)
	}
}
