package main

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var Version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tf-import-gen [flags] address...",
		Short: "Generate terraform import statements",
		Long: strings.TrimSpace(`
Generate terraform import statements to simplify state migrations from one terraform code base to another.

The address argument can be used to filter the instances by resource or module. If
no pattern is given, import statements are generated for all the resources.

The addresses must either be module addresses or absolute resource
addresses, such as:
  aws_instance.example
  module.example
  module.example.module.child
  module.example.aws_instance.example
`),
		Version: Version,
		Example: `
## Generating import statements by module
terraform show -json | tf-import-gen module.example

## Generating import statements by resource
terraform show -json | tf-import-gen aws_instance.example

## Generating import statements by multiple resources
terraform show -json | tf-import-gen aws_instance.example module.example

## Generating import statements for all resources
terraform show -json | tf-import-gen
`,
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
