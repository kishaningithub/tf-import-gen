package main

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg"
	"github.com/spf13/cobra"
	"os"
)

var Version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:     "tf-import-gen",
		Short:   "Generate terraform import statements",
		Long:    "Tool to generate terraform import statements to simplify state migrations from one terraform code base to another",
		Version: Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			address := ""
			if len(os.Args) > 1 {
				address = os.Args[1]
			}
			imports, err := tfimportgen.GenerateImports(os.Stdin, address)
			if err != nil {
				return err
			}
			fmt.Println(imports)
			return nil
		},
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
