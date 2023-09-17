package main

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg"
	"os"
)

func main() {
	address := ""
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	imports, err := tfimportgen.GenerateImports(os.Stdin, address)
	if err != nil {
		panic(err)
	}

	fmt.Println(imports)
}
