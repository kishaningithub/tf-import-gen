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
		return fmt.Sprint(resource.AttributeValues["policy_arn"])
	case "aws_lambda_permission":
		return fmt.Sprint(resource.AttributeValues["statement_id"])
	case "aws_security_group_rule":
		return computeResourceIDForAWSSecurityGroupRole(resource)
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
	cidrBlocks, isCidrBlocksValid := resource.AttributeValues["cidr_blocks"].([]string)
	if isCidrBlocksValid {
		isCidrBlocksValid = len(cidrBlocks) > 0
	}

	resourceID := fmt.Sprintf("%s_%s_%s_%s_%s", securityGroupId, securityGroupType, protocol, fromPort, toPort)
	if isSourceSecurityGroupIdValid {
		return fmt.Sprintf("%s_%s", resourceID, sourceSecurityGroupId)
	}
	if isCidrBlocksValid {
		return fmt.Sprintf("%s_%s", resourceID, strings.Join(cidrBlocks, "_"))
	}
	return resourceID
}
