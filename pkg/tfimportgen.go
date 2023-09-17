package tfimportgen

import (
	"github.com/kishaningithub/tf-import-gen/pkg/internal"
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"io"
)

func GenerateImports(stateJsonReader io.Reader, address string) (internal.TerraformImports, error) {
	resources, err := parser.NewTerraformStateJsonParser(stateJsonReader).Parse()
	if err != nil {
		return nil, err
	}

	resources = resources.FilterByAddress(address)

	var imports internal.TerraformImports
	for _, resource := range resources {
		terraformImport := internal.NewTerraformImport(resource)
		imports = append(imports, terraformImport)
	}

	return imports, nil
}
