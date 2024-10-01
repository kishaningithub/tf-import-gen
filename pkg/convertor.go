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
		"aws_lakeformation_permissions",
		"aws_iam_policy_attachment",
		"aws_acm_certificate_validation",
		"aws_ami_copy",
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
	v := func(name string) string {
		return fmt.Sprint(resource.AttributeValues[name])
	}
	switch resource.Type {
	case "aws_iam_role_policy_attachment":
		return fmt.Sprintf("%s/%s", v("role"), v("policy_arn"))
	case "aws_cloudwatch_event_target":
		return fmt.Sprintf("%s/%s", v("rule"), v("target_id"))
	case "aws_lambda_permission":
		return fmt.Sprintf("%s/%s", v("function_name"), v("statement_id"))
	case "aws_security_group_rule":
		return computeResourceIDForAWSSecurityGroupRole(resource)
	case "aws_api_gateway_resource", "aws_api_gateway_deployment":
		return fmt.Sprintf("%s/%s", v("rest_api_id"), v("id"))
	case "aws_api_gateway_stage":
		return fmt.Sprintf("%s/%s", v("rest_api_id"), v("stage_name"))
	case "aws_api_gateway_method_settings":
		return fmt.Sprintf("%s/%s/%s", v("rest_api_id"), v("stage_name"), v("method_path"))
	case "aws_api_gateway_method", "aws_api_gateway_integration":
		return fmt.Sprintf("%s/%s/%s", v("rest_api_id"), v("resource_id"), v("http_method"))
	case "aws_route_table_association":
		return fmt.Sprintf("%s/%s", v("subnet_id"), v("route_table_id"))
	case "aws_iam_user_policy_attachment":
		return fmt.Sprintf("%s/%s", v("user"), v("policy_arn"))
	case "aws_emr_instance_group":
		return fmt.Sprintf("%s/%s", v("cluster_id"), v("id"))
	case "aws_backup_selection":
		return fmt.Sprintf("%s|%s", v("plan_id"), v("id"))
	case "aws_vpc_endpoint_route_table_association":
		return fmt.Sprintf("%s/%s", v("vpc_endpoint_id"), v("route_table_id"))
	case "aws_cognito_user_pool_client":
		return fmt.Sprintf("%s/%s", v("user_pool_id"), v("id"))
	case "aws_ecs_cluster":
		return v("name")
	case "aws_ecs_task_definition":
		return v("arn")
	case "aws_wafv2_web_acl":
		return fmt.Sprintf("%s/%s/%s", v("id"), v("name"), v("scope"))
	case "aws_autoscaling_schedule":
		return fmt.Sprintf("%s/%s", v("autoscaling_group_name"), v("scheduled_action_name"))
	case "aws_appautoscaling_target":
		return fmt.Sprintf("%s/%s/%s", v("service_namespace"), v("resource_id"), v("scalable_dimension"))
	case "aws_appautoscaling_policy":
		return fmt.Sprintf("%s/%s/%s/%s", v("service_namespace"), v("resource_id"), v("scalable_dimension"), v("name"))
	case "aws_ecs_service":
		return fmt.Sprintf("%s/%s", getEcsClusterNameFromARN(v("cluster")), v("name"))
	default:
		return v("id")
	}
}

func getEcsClusterNameFromARN(arn string) string {
	if parts := strings.Split(arn, "/"); len(parts) == 2 {
		return parts[1]
	}
	return ""
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
	prefixListIds, isPrefixListIdsValid := resource.AttributeValues["prefix_list_ids"].([]any)
	if isPrefixListIdsValid {
		isPrefixListIdsValid = len(prefixListIds) > 0
	}
	resourceID := fmt.Sprintf("%s_%s_%s_%s_%s", securityGroupId, securityGroupType, protocol, fromPort, toPort)
	if isSourceSecurityGroupIdValid {
		return fmt.Sprintf("%s_%s", resourceID, sourceSecurityGroupId)
	}
	if isCidrBlocksValid {
		return fmt.Sprintf("%s_%s", resourceID, strings.Join(convertToStrings(cidrBlocks), "_"))
	}
	if isPrefixListIdsValid {
		return fmt.Sprintf("%s_%s", resourceID, strings.Join(convertToStrings(prefixListIds), "_"))
	}
	return resourceID
}

func convertToStrings(source []any) []string {
	var result []string
	for _, element := range source {
		result = append(result, fmt.Sprint(element))
	}
	return result
}
