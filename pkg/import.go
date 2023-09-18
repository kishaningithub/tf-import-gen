package tfimportgen

import (
	"fmt"
	"strings"
)

var _ fmt.Stringer = TerraformImport{}

type TerraformImport struct {
	ResourceAddress string
	ResourceID      string
}

func (terraformImport TerraformImport) String() string {
	importTemplate := `import {
  to = %s
  id = "%s"
}`
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

func (terraformImport TerraformImport) ToImportBlock() string {
	importTemplate := `import {
  to = %s
  id = "%s"
}`
	return fmt.Sprintf(importTemplate, terraformImport.ResourceAddress, terraformImport.ResourceID)
}
