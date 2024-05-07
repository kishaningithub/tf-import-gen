package tfimportgen

import (
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"io"
)

func GenerateImports(stateJsonReader io.Reader, address ...string) (TerraformImports, error) {
	resources, err := parser.NewTerraformStateJsonParser(stateJsonReader).Parse()
	if err != nil {
		return nil, err
	}

	resources = resources.FilterByAddresses(address...)

	var imports TerraformImports
	for _, resource := range resources {
		terraformImport := computeTerraformImportForResource(resource)
		imports = append(imports, terraformImport)
	}

	return imports, nil
}
