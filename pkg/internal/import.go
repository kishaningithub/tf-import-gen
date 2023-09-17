package internal

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"strings"
)

var _ fmt.Stringer = TerraformImport{}

type TerraformImport struct {
	ResourceAddress string
	ResourceID      string
}

func NewTerraformImport(resource parser.TerraformResource) TerraformImport {
	return TerraformImport{
		ResourceAddress: resource.Address,
		ResourceID:      computeResourceID(resource),
	}
}

func computeResourceID(resource parser.TerraformResource) string {
	switch resource.Type {
	case "aws_iam_role_policy_attachment":
		return fmt.Sprintf("%s", resource.AttributeValues["policy_arn"])
	case "aws_lambda_permission":
		return fmt.Sprintf("%s", resource.AttributeValues["statement_id"])
	default:
		return fmt.Sprintf("%s", resource.AttributeValues["id"])
	}
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
