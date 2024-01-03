package tfimportgen

import (
	"fmt"
	"strings"
)

var _ fmt.Stringer = TerraformImport{}

type TerraformImport struct {
	ResourceAddress string
	ResourceID      string
	SupportsImport  bool
}

func (terraformImport TerraformImport) String() string {
	importTemplate := `import {
  to = %s
  id = "%s"
}`
	if !terraformImport.SupportsImport {
		importTemplate = `# resource "%s" with identifier "%s" does not support import operation. Kindly refer resource documentation for more info.`
	}

	return fmt.Sprintln(fmt.Sprintf(importTemplate, terraformImport.ResourceAddress, terraformImport.ResourceID))
}

var _ fmt.Stringer = (TerraformImports)(nil)

type TerraformImports []TerraformImport

func (terraformImports TerraformImports) String() string {
	var terraformImportsStr strings.Builder
	for _, terraformImport := range terraformImports {
		terraformImportsStr.WriteString(fmt.Sprintln(terraformImport))
	}
	return terraformImportsStr.String()
}
