package tfimportgen

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"slices"
	"strings"
)

func computeTerraformImportForResource(resource parser.TerraformResource) TerraformImport {
	resourcesWhichDoNotSupportImport := []string{
		"aws_alb_target_group_attachment",
		"aws_lb_target_group_attachment",
		"aws_lakeformation_data_lake_settings",
	}
	if slices.Contains(resourcesWhichDoNotSupportImport, resource.Type) {
		return TerraformImport{
			SupportsImport:  false,
			ResourceAddress: resource.Address,
			ResourceID:      computeResourceID(resource),
		}
	}
	return TerraformImport{
		SupportsImport:  true,
		ResourceAddress: resource.Address,
		ResourceID:      computeResourceID(resource),
	}
}

func computeResourceID(resource parser.TerraformResource) string {
	getValue := func(name string) string {
		return fmt.Sprint(resource.AttributeValues[name])
	}
	switch resource.Type {
	case "aws_iam_role_policy_attachment":
		return fmt.Sprintf("%s/%s", getValue("role"), getValue("policy_arn"))
	case "aws_cloudwatch_event_target":
		return fmt.Sprintf("%s/%s", getValue("rule"), getValue("target_id"))
	case "aws_lambda_permission":
		return fmt.Sprintf("%s/%s", getValue("function_name"), getValue("statement_id"))
	case "aws_security_group_rule":
		return computeResourceIDForAWSSecurityGroupRole(resource)
	case "aws_api_gateway_resource", "aws_api_gateway_deployment":
		return fmt.Sprintf("%s/%s", getValue("rest_api_id"), getValue("id"))
	case "aws_api_gateway_stage":
		return fmt.Sprintf("%s/%s", getValue("rest_api_id"), getValue("stage_name"))
	case "aws_api_gateway_method_settings":
		return fmt.Sprintf("%s/%s/%s", getValue("rest_api_id"), getValue("stage_name"), getValue("method_path"))
	case "aws_api_gateway_method", "aws_api_gateway_integration":
		return fmt.Sprintf("%s/%s/%s", getValue("rest_api_id"), getValue("resource_id"), getValue("http_method"))
	case "aws_route_table_association":
		return fmt.Sprintf("%s/%s", getValue("subnet_id"), getValue("route_table_id"))
	case "aws_iam_user_policy_attachment":
		return fmt.Sprintf("%s/%s", getValue("user"), getValue("policy_arn"))
	default:
		return getValue("id")
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
