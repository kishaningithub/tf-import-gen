package tfimportgen

import (
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"io"
)

func GenerateImports(stateJsonReader io.Reader, addresses []string) (TerraformImports, error) {
	resources, err := parser.NewTerraformStateJsonParser(stateJsonReader).Parse()
	if err != nil {
		return nil, err
	}

	if addresses != nil {
		resources = resources.FilterByAddresses(addresses)
	}

	var imports TerraformImports
	for _, resource := range resources {
		terraformImport := computeTerraformImportForResource(resource)
		imports = append(imports, terraformImport)
	}

	return imports, nil
}
