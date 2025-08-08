package tfimportgen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
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
		"google_storage_default_object",
		"google_storage_default_object_acl",
		"google_storage_bucket_acl",
		"google_project_service_identity",
		"google_project_default_service_accounts",
		"google_compute_instance_from_template",
		"google_pubsub_subscription_iam_member",
		"google_pubsub_subscription_iam_binding",
		"google_compute_instance_template",
		"google_iap_tunnel_instance_iam_binding",
		"local_file",
		"tailscale_tailnet_key",
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
	// aws resources
	case "aws_iam_role_policy_attachment":
		return fmt.Sprintf("%s/%s", v("role"), v("policy_arn"))
	case "aws_cloudwatch_event_target":
		return fmt.Sprintf("%s/%s", v("rule"), v("target_id"))
	case "aws_lambda_permission":
		return fmt.Sprintf("%s/%s", v("function_name"), v("statement_id"))
	case "aws_security_group_rule":
		return computeResourceIDForAWSSecurityGroupRole(resource)
	case "aws_network_acl_rule":
		return computeResourceIdForAWSNetworkACLRule(resource)
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
	case "aws_vpc_endpoint_subnet_association":
		return fmt.Sprintf("%s/%s", v("vpc_endpoint_id"), v("subnet_id"))
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
	case "aws_cloudwatch_log_stream":
		return fmt.Sprintf("%s:%s", v("log_group_name"), v("name"))
	case "aws_route":
		if resource.AttributeValues["destination_prefix_list_id"] != "" {
			return fmt.Sprintf("%s_%s", v("route_table_id"), v("destination_prefix_list_id"))
		} else if resource.AttributeValues["destination_cidr_block"] != "" {
			return fmt.Sprintf("%s_%s", v("route_table_id"), v("destination_cidr_block"))
		} else {
			return fmt.Sprintf("%s_%s", v("route_table_id"), v("destination_ipv6_cidr_block"))
		}
	// gcp resources
	case "google_bigquery_dataset_iam_member":
		return fmt.Sprintf("projects/%s/datasets/%s %s %s", v("project"), v("dataset_id"), v("role"), v("member"))
	case "google_bigquery_table_iam_member":
		return fmt.Sprintf("%s %s %s", v("table_id"), v("role"), v("member"))
	case "google_service_account_iam_member":
		return fmt.Sprintf("%s %s %s", v("service_account_id"), v("role"), v("member"))
	case "google_service_account_iam_binding":
		return fmt.Sprintf("%s %s", v("service_account_id"), v("role"))
	case "google_privateca_ca_pool_iam_member":
		conditions, ok := resource.AttributeValues["condition"].([]any)
		if ok && len(conditions) > 0 {
			condition := conditions[0].(map[string]any)
			return fmt.Sprintf("%s %s %s %s", v("ca_pool"), v("role"), v("member"), condition["title"])
		}
		return fmt.Sprintf("%s %s %s", v("ca_pool"), v("role"), v("member"))
	case "google_privateca_certificate_template_iam_member":
		return fmt.Sprintf("%s %s %s", v("certificate_template"), v("role"), v("member"))
	case "google_cloud_run_service_iam_binding":
		return fmt.Sprintf("%s %s", v("service"), v("role"))
	case "google_kms_crypto_key_iam_binding":
		return fmt.Sprintf("%s %s", v("crypto_key_id"), v("role"))
	case "google_kms_crypto_key_iam_member":
		return fmt.Sprintf("%s %s %s", v("crypto_key_id"), v("role"), v("member"))
	case "google_organization_iam_member":
		return fmt.Sprintf("%s %s %s", v("org_id"), v("role"), v("member"))
	case "google_project_iam_member":
		conditions, ok := resource.AttributeValues["condition"].([]any)
		if ok && len(conditions) > 0 {
			condition := conditions[0].(map[string]any)
			return fmt.Sprintf("%s %s %s %s", v("project"), v("role"), v("member"), condition["title"])
		}
		return fmt.Sprintf("%s %s %s", v("project"), v("role"), v("member"))
	case "google_project_iam_binding":
		return fmt.Sprintf("%s %s", v("project"), v("role"))
	case "google_project_iam_custom_role":
		return fmt.Sprintf("%s %s", v("project"), v("id"))
	case "google_sql_database_instance":
		return fmt.Sprintf("projects/%s/instances/%s", v("project"), v("name"))
	case "google_sql_user":
		return fmt.Sprintf("%s/%s/%s", v("project"), v("instance"), v("name"))
	case "google_iap_tunnel_instance_iam_binding":
		return fmt.Sprintf("%s %s", v("instance"), v("role"))
	case "google_secret_manager_secret_iam_binding":
		return fmt.Sprintf("%s %s", v("secret_id"), v("role"))
	case "google_secret_manager_secret_iam_member":
		return fmt.Sprintf("%s %s %s", v("secret_id"), v("role"), v("member"))
	case "google_secret_manager_secret":
		return v("name")
	case "google_storage_bucket_iam_member":
		return fmt.Sprintf("%s %s %s", v("bucket"), v("role"), v("member"))
	case "google_storage_bucket_iam_binding":
		return fmt.Sprintf("%s %s", v("bucket"), v("role"))
	case "google_tags_tag_key_iam_member":
		return fmt.Sprintf("%s %s %s", v("tag_key"), v("role"), v("member"))
	case "google_compute_subnetwork_iam_binding":
		return fmt.Sprintf("%s %s", v("subnetwork"), v("role"))
	case "google_pubsub_topic_iam_binding":
		return fmt.Sprintf("%s %s", v("topic"), v("role"))
	case "google_pubsub_topic_iam_member":
		conditions, ok := resource.AttributeValues["condition"].([]any)
		if ok && len(conditions) > 0 {
			condition := conditions[0].(map[string]any)
			return fmt.Sprintf("%s %s %s %s", v("topic"), v("role"), v("member"), condition["title"])
		}
		return fmt.Sprintf("%s %s %s", v("topic"), v("role"), v("member"))
	case "google_resource_manager_lien":
		return fmt.Sprintf("%s/%s", strings.ReplaceAll(v("parent"), "projects/", ""), v("name"))
	case "google_monitoring_uptime_check_config":
		return fmt.Sprintf("%s %s", v("project"), v("id"))
	case "google_monitoring_alert_policy":
		return fmt.Sprintf("%s %s", v("project"), v("name"))
	case "google_monitoring_notification_channel":
		return v("name")
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

func computeResourceIdForAWSNetworkACLRule(resource parser.TerraformResource) string {
	networkACLId := fmt.Sprint(resource.AttributeValues["network_acl_id"])
	ruleNumber := fmt.Sprint(resource.AttributeValues["rule_number"])
	protocol := fmt.Sprint(resource.AttributeValues["protocol"])
	egress := fmt.Sprint(resource.AttributeValues["egress"])
	return fmt.Sprintf("%s:%s:%s:%s", networkACLId, ruleNumber, protocol, egress)
}

func convertToStrings(source []any) []string {
	var result []string
	for _, element := range source {
		result = append(result, fmt.Sprint(element))
	}
	return result
}
