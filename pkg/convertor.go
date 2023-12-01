package tfimportgen

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"strings"
)

func computeTerraformImportForResource(resource parser.TerraformResource) TerraformImport {
	return TerraformImport{
		ResourceAddress: resource.Address,
		ResourceID:      computeResourceID(resource),
	}
}

func computeResourceID(resource parser.TerraformResource) string {
	switch resource.Type {
	case "aws_iam_role_policy_attachment":
		role := fmt.Sprint(resource.AttributeValues["role"])
		policyArn := fmt.Sprint(resource.AttributeValues["policy_arn"])
		return fmt.Sprintf("%s/%s", role, policyArn)
	case "aws_lambda_permission":
		functionName := fmt.Sprint(resource.AttributeValues["function_name"])
		statementId := fmt.Sprint(resource.AttributeValues["statement_id"])
		return fmt.Sprintf("%s/%s", functionName, statementId)
	case "aws_security_group_rule":
		return computeResourceIDForAWSSecurityGroupRole(resource)
	case "aws_api_gateway_resource":
		restApiId := fmt.Sprint(resource.AttributeValues["rest_api_id"])
		resourceId := fmt.Sprint(resource.AttributeValues["id"])
		return fmt.Sprintf("%s/%s", restApiId, resourceId)
	default:
		return fmt.Sprint(resource.AttributeValues["id"])
	}
}

func computeResourceIDForAWSSecurityGroupRole(resource parser.TerraformResource) string {
	// Required Fields
	securityGroupId := fmt.Sprint(resource.AttributeValues["security_group_id"])
	securityGroupType := fmt.Sprint(resource.AttributeValues["type"])
	protocol := fmt.Sprint(resource.AttributeValues["protocol"])
	fromPort := fmt.Sprint(resource.AttributeValues["from_port"])
	toPort := fmt.Sprint(resource.AttributeValues["to_port"])

	// Optional Fields
	sourceSecurityGroupId, isSourceSecurityGroupIdValid := resource.AttributeValues["source_security_group_id"].(string)
	if isSourceSecurityGroupIdValid {
		isSourceSecurityGroupIdValid = len(sourceSecurityGroupId) > 0
	}
	cidrBlocks, isCidrBlocksValid := resource.AttributeValues["cidr_blocks"].([]any)
	if isCidrBlocksValid {
		isCidrBlocksValid = len(cidrBlocks) > 0
	}
	resourceID := fmt.Sprintf("%s_%s_%s_%s_%s", securityGroupId, securityGroupType, protocol, fromPort, toPort)
	if isSourceSecurityGroupIdValid {
		return fmt.Sprintf("%s_%s", resourceID, sourceSecurityGroupId)
	}
	if isCidrBlocksValid {
		var cidrStringBlocks []string
		for _, cidrBlock := range cidrBlocks {
			cidrStringBlocks = append(cidrStringBlocks, fmt.Sprint(cidrBlock))
		}
		return fmt.Sprintf("%s_%s", resourceID, strings.Join(cidrStringBlocks, "_"))
	}
	return resourceID
}
