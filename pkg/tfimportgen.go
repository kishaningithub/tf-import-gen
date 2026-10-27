package tfimportgen

import (
	"io"

	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"github.com/samber/lo"
)

func GenerateImports(stateJsonReader io.Reader, addresses []string) (TerraformImports, error) {
	resources, err := parser.NewTerraformStateJsonParser(stateJsonReader).Parse()
	if err != nil {
		return nil, err
	}

	if addresses != nil {
		resources = resources.FilterByAddresses(addresses)
	}

	imports := lo.Map(resources, func(resource parser.TerraformResource, _ int) TerraformImport {
		return computeTerraformImportForResource(resource)
	})

	return imports, nil
}
