package tfimportgen

import (
	"testing"

	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"github.com/stretchr/testify/require"
)

func Test_ComputeTerraformImportForResource(t *testing.T) {
	tests := []struct {
		name              string
		terraformResource parser.TerraformResource
		expected          TerraformImport
	}{
		{
			name: "For aws_iam_role_policy_attachment",
			terraformResource: parser.TerraformResource{
				Address: "aws_iam_role_policy_attachment.test",
				Type:    "aws_iam_role_policy_attachment",
				AttributeValues: map[string]any{
					"role":       "test-role",
					"policy_arn": "test-policy-arn",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_iam_role_policy_attachment.test",
				ResourceID:      "test-role/test-policy-arn",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_lambda_permission",
			terraformResource: parser.TerraformResource{
				Address: "aws_lambda_permission.test",
				Type:    "aws_lambda_permission",
				AttributeValues: map[string]any{
					"statement_id":  "test-statement-id",
					"function_name": "test-function-name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_lambda_permission.test",
				ResourceID:      "test-function-name/test-statement-id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_security_group_rule with source_security_group_id",
			terraformResource: parser.TerraformResource{
				Address: "aws_security_group_rule.test",
				Type:    "aws_security_group_rule",
				AttributeValues: map[string]any{
					"security_group_id":        "security-group-id",
					"type":                     "type",
					"protocol":                 "protocol",
					"from_port":                1234,
					"to_port":                  5678,
					"source_security_group_id": "source-security-group-id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_security_group_rule.test",
				ResourceID:      "security-group-id_type_protocol_1234_5678_source-security-group-id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_security_group_rule with cidr_blocks",
			terraformResource: parser.TerraformResource{
				Address: "aws_security_group_rule.test",
				Type:    "aws_security_group_rule",
				AttributeValues: map[string]any{
					"security_group_id": "security-group-id",
					"type":              "type",
					"protocol":          "protocol",
					"from_port":         1234,
					"to_port":           5678,
					"cidr_blocks":       []any{"cidr-block-1", "cidr-block-2"},
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_security_group_rule.test",
				ResourceID:      "security-group-id_type_protocol_1234_5678_cidr-block-1_cidr-block-2",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_security_group_rule with prefix_list_ids",
			terraformResource: parser.TerraformResource{
				Address: "aws_security_group_rule.test",
				Type:    "aws_security_group_rule",
				AttributeValues: map[string]any{
					"security_group_id": "security-group-id",
					"type":              "type",
					"protocol":          "protocol",
					"from_port":         1234,
					"to_port":           5678,
					"prefix_list_ids":   []any{"prefix-list-1", "prefix-list-2"},
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_security_group_rule.test",
				ResourceID:      "security-group-id_type_protocol_1234_5678_prefix-list-1_prefix-list-2",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_resource",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_resource.test",
				Type:    "aws_api_gateway_resource",
				AttributeValues: map[string]any{
					"id":          "id",
					"rest_api_id": "rest_api_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_resource.test",
				ResourceID:      "rest_api_id/id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_deployment",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_deployment.test",
				Type:    "aws_api_gateway_deployment",
				AttributeValues: map[string]any{
					"id":          "id",
					"rest_api_id": "rest_api_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_deployment.test",
				ResourceID:      "rest_api_id/id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_stage",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_stage.test",
				Type:    "aws_api_gateway_stage",
				AttributeValues: map[string]any{
					"rest_api_id": "rest_api_id",
					"stage_name":  "stage_name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_stage.test",
				ResourceID:      "rest_api_id/stage_name",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_method_settings",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_method_settings.test",
				Type:    "aws_api_gateway_method_settings",
				AttributeValues: map[string]any{
					"rest_api_id": "rest_api_id",
					"stage_name":  "stage_name",
					"method_path": "method_path",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_method_settings.test",
				ResourceID:      "rest_api_id/stage_name/method_path",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_cloudwatch_event_target",
			terraformResource: parser.TerraformResource{
				Address: "aws_cloudwatch_event_target.test",
				Type:    "aws_cloudwatch_event_target",
				AttributeValues: map[string]any{
					"rule":      "rule",
					"target_id": "target_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_cloudwatch_event_target.test",
				ResourceID:      "rule/target_id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_method",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_method.test",
				Type:    "aws_api_gateway_method",
				AttributeValues: map[string]any{
					"http_method": "http_method",
					"resource_id": "resource_id",
					"rest_api_id": "rest_api_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_method.test",
				ResourceID:      "rest_api_id/resource_id/http_method",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_api_gateway_integration",
			terraformResource: parser.TerraformResource{
				Address: "aws_api_gateway_integration.test",
				Type:    "aws_api_gateway_integration",
				AttributeValues: map[string]any{
					"http_method": "http_method",
					"resource_id": "resource_id",
					"rest_api_id": "rest_api_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_api_gateway_integration.test",
				ResourceID:      "rest_api_id/resource_id/http_method",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_route_table_association",
			terraformResource: parser.TerraformResource{
				Address: "aws_route_table_association.test",
				Type:    "aws_route_table_association",
				AttributeValues: map[string]any{
					"subnet_id":      "subnet_id",
					"route_table_id": "route_table_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_route_table_association.test",
				ResourceID:      "subnet_id/route_table_id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_iam_user_policy_attachment",
			terraformResource: parser.TerraformResource{
				Address: "aws_iam_user_policy_attachment.test",
				Type:    "aws_iam_user_policy_attachment",
				AttributeValues: map[string]any{
					"policy_arn": "policy_arn",
					"user":       "user",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_iam_user_policy_attachment.test",
				ResourceID:      "user/policy_arn",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_emr_instance_group",
			terraformResource: parser.TerraformResource{
				Address: "aws_emr_instance_group.test",
				Type:    "aws_emr_instance_group",
				AttributeValues: map[string]any{
					"id":         "id",
					"cluster_id": "cluster_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_emr_instance_group.test",
				ResourceID:      "cluster_id/id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_backup_selection",
			terraformResource: parser.TerraformResource{
				Address: "aws_backup_selection.test",
				Type:    "aws_backup_selection",
				AttributeValues: map[string]any{
					"id":      "id",
					"plan_id": "plan_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_backup_selection.test",
				ResourceID:      "plan_id|id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_vpc_endpoint_route_table_association",
			terraformResource: parser.TerraformResource{
				Address: "aws_vpc_endpoint_route_table_association.test",
				Type:    "aws_vpc_endpoint_route_table_association",
				AttributeValues: map[string]any{
					"vpc_endpoint_id": "vpc_endpoint_id",
					"route_table_id":  "route_table_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_vpc_endpoint_route_table_association.test",
				ResourceID:      "vpc_endpoint_id/route_table_id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_vpc_endpoint_subnet_association",
			terraformResource: parser.TerraformResource{
				Address: "aws_vpc_endpoint_subnet_association.test",
				Type:    "aws_vpc_endpoint_subnet_association",
				AttributeValues: map[string]any{
					"vpc_endpoint_id": "vpce-12345",
					"subnet_id":       "subnet-67890",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_vpc_endpoint_subnet_association.test",
				ResourceID:      "vpce-12345/subnet-67890",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_network_acl_rule",
			terraformResource: parser.TerraformResource{
				Address: "aws_network_acl_rule.test",
				Type:    "aws_network_acl_rule",
				AttributeValues: map[string]any{
					"network_acl_id": "acl-12345",
					"rule_number":    100,
					"protocol":       "6",
					"egress":         false,
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_network_acl_rule.test",
				ResourceID:      "acl-12345:100:6:false",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_cognito_user_pool_client",
			terraformResource: parser.TerraformResource{
				Address: "aws_cognito_user_pool_client.test",
				Type:    "aws_cognito_user_pool_client",
				AttributeValues: map[string]any{
					"user_pool_id": "user_pool_id",
					"id":           "id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_cognito_user_pool_client.test",
				ResourceID:      "user_pool_id/id",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_ecs_cluster",
			terraformResource: parser.TerraformResource{
				Address: "aws_ecs_cluster.test",
				Type:    "aws_ecs_cluster",
				AttributeValues: map[string]any{
					"name": "name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_ecs_cluster.test",
				ResourceID:      "name",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_ecs_task_definition",
			terraformResource: parser.TerraformResource{
				Address: "aws_ecs_task_definition.test",
				Type:    "aws_ecs_task_definition",
				AttributeValues: map[string]any{
					"arn": "arn",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_ecs_task_definition.test",
				ResourceID:      "arn",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_wafv2_web_acl",
			terraformResource: parser.TerraformResource{
				Address: "aws_wafv2_web_acl.test",
				Type:    "aws_wafv2_web_acl",
				AttributeValues: map[string]any{
					"id":    "id",
					"name":  "name",
					"scope": "scope",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_wafv2_web_acl.test",
				ResourceID:      "id/name/scope",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_autoscaling_schedule",
			terraformResource: parser.TerraformResource{
				Address: "aws_autoscaling_schedule.test",
				Type:    "aws_autoscaling_schedule",
				AttributeValues: map[string]any{
					"autoscaling_group_name": "autoscaling_group_name",
					"scheduled_action_name":  "scheduled_action_name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_autoscaling_schedule.test",
				ResourceID:      "autoscaling_group_name/scheduled_action_name",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_appautoscaling_target",
			terraformResource: parser.TerraformResource{
				Address: "aws_appautoscaling_target.test",
				Type:    "aws_appautoscaling_target",
				AttributeValues: map[string]any{
					"service_namespace":  "service_namespace",
					"resource_id":        "resource_id",
					"scalable_dimension": "scalable_dimension",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_appautoscaling_target.test",
				ResourceID:      "service_namespace/resource_id/scalable_dimension",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_appautoscaling_policy",
			terraformResource: parser.TerraformResource{
				Address: "aws_appautoscaling_policy.test",
				Type:    "aws_appautoscaling_policy",
				AttributeValues: map[string]any{
					"service_namespace":  "service_namespace",
					"resource_id":        "resource_id",
					"scalable_dimension": "scalable_dimension",
					"name":               "policy_name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_appautoscaling_policy.test",
				ResourceID:      "service_namespace/resource_id/scalable_dimension/policy_name",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_ecs_service",
			terraformResource: parser.TerraformResource{
				Address: "aws_ecs_service.test",
				Type:    "aws_ecs_service",
				AttributeValues: map[string]any{
					"name":    "service_name",
					"cluster": "arn:aws:ecs:us-west-2:0123456789:cluster/cluster_name",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_ecs_service.test",
				ResourceID:      "cluster_name/service_name",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_cloudwatch_log_stream",
			terraformResource: parser.TerraformResource{
				Address: "aws_cloudwatch_log_stream.test",
				Type:    "aws_cloudwatch_log_stream",
				AttributeValues: map[string]any{
					"log_group_name": "test_log_group",
					"name":           "TestLogStream123",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_cloudwatch_log_stream.test",
				ResourceID:      "test_log_group:TestLogStream123",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_route_with_destination_prefix_list_id",
			terraformResource: parser.TerraformResource{
				Address: "aws_route.test",
				Type:    "aws_route",
				AttributeValues: map[string]any{
					"route_table_id":              "rtb-656C65616E6F72",
					"destination_prefix_list_id":  "pl-12df45133",
					"destination_cidr_block":      "",
					"destination_ipv6_cidr_block": "",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_route.test",
				ResourceID:      "rtb-656C65616E6F72_pl-12df45133",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_route_with_destination_cidr_block",
			terraformResource: parser.TerraformResource{
				Address: "aws_route.test",
				Type:    "aws_route",
				AttributeValues: map[string]any{
					"route_table_id":              "rtb-656C65616E6F72",
					"destination_cidr_block":      "10.42.0.0/16",
					"destination_ipv6_cidr_block": "",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_route.test",
				ResourceID:      "rtb-656C65616E6F72_10.42.0.0/16",
				SupportsImport:  true,
			},
		},
		{
			name: "For aws_route_with_destination_ipv6_cidr_block",
			terraformResource: parser.TerraformResource{
				Address: "aws_route.test",
				Type:    "aws_route",
				AttributeValues: map[string]any{
					"route_table_id":              "rtb-656C65616E6F72",
					"destination_cidr_block":      "",
					"destination_ipv6_cidr_block": "2620:0:2d0:200::8/1",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "aws_route.test",
				ResourceID:      "rtb-656C65616E6F72_2620:0:2d0:200::8/1",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_monitoring_alert_policy",
			terraformResource: parser.TerraformResource{
				Address: "google_monitoring_alert_policy.test",
				Type:    "google_monitoring_alert_policy",
				AttributeValues: map[string]any{
					"name":    "projects/project/alertPolicies/123456789",
					"project": "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_monitoring_alert_policy.test",
				ResourceID:      "project projects/project/alertPolicies/123456789",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_monitoring_notification_channel",
			terraformResource: parser.TerraformResource{
				Address: "google_monitoring_notification_channel.test",
				Type:    "google_monitoring_notification_channel",
				AttributeValues: map[string]any{
					"display_name": "Team Email Group",
					"id":           "project projects/project/notificationChannels/12345678901234567",
					"name":         "projects/project/notificationChannels/12345678901234567",
					"project":      "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_monitoring_notification_channel.test",
				ResourceID:      "projects/project/notificationChannels/12345678901234567",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_service_account_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_service_account_iam_binding.test",
				Type:    "google_service_account_iam_binding",
				AttributeValues: map[string]any{
					"id":                 "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com/roles/iam.serviceAccountUser",
					"role":               "roles/iam.serviceAccountUser",
					"service_account_id": "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_service_account_iam_binding.test",
				ResourceID:      "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com roles/iam.serviceAccountUser",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_service_account_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_service_account_iam_member.test",
				Type:    "google_service_account_iam_member",
				AttributeValues: map[string]any{
					"id":                 "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com/roles/iam.serviceAccountUser/user:user.name@email.com",
					"member":             "user:user.name@email.com",
					"role":               "roles/iam.serviceAccountUser",
					"service_account_id": "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_service_account_iam_member.test",
				ResourceID:      "projects/project/serviceAccounts/service@project.iam.gserviceaccount.com roles/iam.serviceAccountUser user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_tags_tag_key_iam_member.tailscale_services_tag_user",
			terraformResource: parser.TerraformResource{
				Address: "google_tags_tag_key_iam_member.tailscale_services_tag_user",
				Type:    "google_tags_tag_key_iam_member",
				AttributeValues: map[string]any{
					"id":      "tagKeys/123456789012345/projects/project/roles/customTailscaleRole/serviceAccount:service@project.iam.gserviceaccount.com",
					"member":  "serviceAccount:service@project.iam.gserviceaccount.com",
					"role":    "projects/project/roles/customTailscaleRole",
					"tag_key": "tagKeys/123456789012345",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_tags_tag_key_iam_member.tailscale_services_tag_user",
				ResourceID:      "tagKeys/123456789012345 projects/project/roles/customTailscaleRole serviceAccount:service@project.iam.gserviceaccount.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_privateca_ca_pool_iam_member.pool",
			terraformResource: parser.TerraformResource{
				Address: "google_privateca_ca_pool_iam_member.pool",
				Type:    "google_privateca_ca_pool_iam_member",
				AttributeValues: map[string]any{
					"ca_pool":  "projects/project/locations/us-north1/caPools/service",
					"id":       "projects/project/locations/us-north1/caPools/service/roles/privateca.poolReader/user:user.name@email.com",
					"member":   "user:user.name@email.com",
					"location": "us-north1",
					"project":  "project",
					"role":     "roles/privateca.poolReader",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_privateca_ca_pool_iam_member.pool",
				ResourceID:      "projects/project/locations/us-north1/caPools/service roles/privateca.poolReader user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_privateca_certificate_template_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_privateca_certificate_template_iam_member.test",
				Type:    "google_privateca_certificate_template_iam_member",
				AttributeValues: map[string]any{
					"certificate_template": "projects/project/locations/us-north1/certificateTemplates/test",
					"id":                   "projects/project/locations/us-north1/certificateTemplates/test/roles/privateca.templateUser/user:user.name@email.com",
					"location":             "us-north1",
					"member":               "user:user.name@email.com",
					"project":              "project",
					"role":                 "roles/privateca.templateUser",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_privateca_certificate_template_iam_member.test",
				ResourceID:      "projects/project/locations/us-north1/certificateTemplates/test roles/privateca.templateUser user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_project_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_project_iam_binding.test",
				Type:    "google_project_iam_binding",
				AttributeValues: map[string]any{
					"id": "project/projects/project/roles/testuser",
					"members": []string{
						"serviceAccount:service@project.iam.gserviceaccount.com",
						"serviceAccount:test@project.iam.gserviceaccount.com",
					},
					"project": "project",
					"role":    "projects/project/roles/testuser",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_project_iam_binding.test",
				ResourceID:      "project projects/project/roles/testuser",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_project_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_project_iam_member.test",
				Type:    "google_project_iam_member",
				AttributeValues: map[string]any{
					"id":      "project/projects/project/roles/test/serviceAccount:service_account@company.iam.gserviceaccount.com",
					"member":  "serviceAccount:service_account@company.iam.gserviceaccount.com",
					"project": "project",
					"role":    "projects/project/roles/test",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_project_iam_member.test",
				ResourceID:      "project projects/project/roles/test serviceAccount:service_account@company.iam.gserviceaccount.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_sql_database_instance",
			terraformResource: parser.TerraformResource{
				Address: "google_sql_database_instance.test",
				Type:    "google_sql_database_instance",
				AttributeValues: map[string]any{
					"id":      "test-ydfs787sd8f",
					"name":    "test-ydfs787sd8f",
					"project": "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_sql_database_instance.test",
				ResourceID:      "projects/project/instances/test-ydfs787sd8f",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_sql_user",
			terraformResource: parser.TerraformResource{
				Address: "google_sql_user.test",
				Type:    "google_sql_user",
				AttributeValues: map[string]any{
					"id":       "service@project.iam//instance-test-45l3n534jh5n",
					"instance": "instance-test-45l3n534jh5n",
					"name":     "service@project.iam",
					"project":  "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_sql_user.test",
				ResourceID:      "project/instance-test-45l3n534jh5n/service@project.iam",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_iap_tunnel_instance_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_iap_tunnel_instance_iam_binding.test",
				Type:    "google_iap_tunnel_instance_iam_binding",
				AttributeValues: map[string]any{
					"id":       "projects/project/iap_tunnel/zones/us-north1-c/instances/test/roles/iap.tunnelResourceAccessor",
					"instance": "projects/project/iap_tunnel/zones/us-north1-c/instances/test",
					"members": []string{
						"group:team@company.com",
						"user:user.name@email.com",
					},
					"project": "project",
					"role":    "roles/iap.tunnelResourceAccessor",
					"zone":    "us-north1-c",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_iap_tunnel_instance_iam_binding.test",
				ResourceID:      "projects/project/iap_tunnel/zones/us-north1-c/instances/test roles/iap.tunnelResourceAccessor",
				SupportsImport:  false,
			},
		},
		{
			name: "For google_bigquery_dataset_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_bigquery_dataset_iam_member.test",
				Type:    "google_bigquery_dataset_iam_member",
				AttributeValues: map[string]any{
					"dataset_id": "test_123456",
					"etag":       "",
					"id":         "projects/project/datasets/test_123456/roles/bigquery.dataEditor/user:user.name@email.com",
					"member":     "user:user.name@email.com",
					"project":    "project",
					"role":       "roles/bigquery.dataEditor",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_bigquery_dataset_iam_member.test",
				ResourceID:      "projects/project/datasets/test_123456 roles/bigquery.dataEditor user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_bigquery_table_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_bigquery_table_iam_member.test",
				Type:    "google_bigquery_table_iam_member",
				AttributeValues: map[string]any{
					"dataset_id": "test_123456",
					"id":         "projects/project/datasets/test_123456/tables/_retracted_file/roles/bigquery.dataEditor/group:group.name@email.com",
					"member":     "group:group.name@email.com",
					"project":    "project",
					"role":       "roles/bigquery.dataEditor",
					"table_id":   "projects/project/datasets/test_123456/tables/_retracted_file",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_bigquery_table_iam_member.test",
				ResourceID:      "projects/project/datasets/test_123456/tables/_retracted_file roles/bigquery.dataEditor group:group.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_kms_crypto_key_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_kms_crypto_key_iam_member.test",
				Type:    "google_kms_crypto_key_iam_member",
				AttributeValues: map[string]any{
					"crypto_key_id": "projects/project/locations/europe/keyRings/test-data/cryptoKeys/data-converter",
					"id":            "projects/project/locations/europe/keyRings/test-data/cryptoKeys/data-converter/roles/cloudkms.cryptoKeyEncrypterDecrypter/serviceAccount:service@project.iam.gserviceaccount.com",
					"member":        "serviceAccount:service@project.iam.gserviceaccount.com",
					"role":          "roles/cloudkms.cryptoKeyEncrypterDecrypter",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_kms_crypto_key_iam_member.test",
				ResourceID:      "projects/project/locations/europe/keyRings/test-data/cryptoKeys/data-converter roles/cloudkms.cryptoKeyEncrypterDecrypter serviceAccount:service@project.iam.gserviceaccount.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_kms_crypto_key_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_kms_crypto_key_iam_binding.test",
				Type:    "google_kms_crypto_key_iam_binding",
				AttributeValues: map[string]any{
					"crypto_key_id": "projects/project/locations/us-north1/keyRings/gke-backup/cryptoKeys/backup",
					"id":            "projects/project/locations/us-north1/keyRings/gke-backup/cryptoKeys/backup/roles/cloudkms.cryptoKeyDecrypter",
					"members": []string{
						"serviceAccount:service@container-engine-robot.iam.gserviceaccount.com",
						"serviceAccount:service@gcp-sa-gkebackup.iam.gserviceaccount.com",
					},
					"role": "roles/cloudkms.cryptoKeyDecrypter",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_kms_crypto_key_iam_binding.test",
				ResourceID:      "projects/project/locations/us-north1/keyRings/gke-backup/cryptoKeys/backup roles/cloudkms.cryptoKeyDecrypter",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_cloud_run_service_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_cloud_run_service_iam_binding.test",
				Type:    "google_cloud_run_service_iam_binding",
				AttributeValues: map[string]any{
					"id":       "v1/projects/project/locations/us-north1/services/service/roles/run.invoker",
					"location": "us-north1",
					"members": []string{
						"serviceAccount:service@project.iam.gserviceaccount.com",
					},
					"project": "project",
					"role":    "roles/run.invoker",
					"service": "v1/projects/project/locations/us-north1/services/service",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_cloud_run_service_iam_binding.test",
				ResourceID:      "v1/projects/project/locations/us-north1/services/service roles/run.invoker",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_secret_manager_secret_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_secret_manager_secret_iam_member.test",
				Type:    "google_secret_manager_secret_iam_member",
				AttributeValues: map[string]any{
					"id":        "projects/project/secrets/test-developer-client-key/roles/secretmanager.secretAccessor/user:user.name@email.com",
					"member":    "user:user.name@email.com",
					"project":   "project",
					"role":      "roles/secretmanager.secretAccessor",
					"secret_id": "projects/project/secrets/test-developer-client-key",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_secret_manager_secret_iam_member.test",
				ResourceID:      "projects/project/secrets/test-developer-client-key roles/secretmanager.secretAccessor user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_secret_manager_secret_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_secret_manager_secret_iam_binding.test",
				Type:    "google_secret_manager_secret_iam_binding",
				AttributeValues: map[string]any{
					"id": "projects/project/secrets/service-import/roles/secretmanager.secretAccessor",
					"members": []string{
						"serviceAccount:service@project.iam.gserviceaccount.com",
						"serviceAccount:test@project.iam.gserviceaccount.com",
					},
					"project":   "project",
					"role":      "roles/secretmanager.secretAccessor",
					"secret_id": "projects/project/secrets/service-import",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_secret_manager_secret_iam_binding.test",
				ResourceID:      "projects/project/secrets/service-import roles/secretmanager.secretAccessor",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_secret_manager_secret.gcp_connector_secret",
			terraformResource: parser.TerraformResource{
				Address: "google_secret_manager_secret.gcp_connector_secret",
				Type:    "google_secret_manager_secret",
				AttributeValues: map[string]any{
					"id":      "projects/project/secrets/a-secret-key",
					"name":    "projects/1234567890/secrets/a-secret-key",
					"project": "us-shared-prod",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_secret_manager_secret.gcp_connector_secret",
				ResourceID:      "projects/1234567890/secrets/a-secret-key",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_storage_bucket_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_storage_bucket_iam_member.test",
				Type:    "google_storage_bucket_iam_member",
				AttributeValues: map[string]any{
					"bucket": "b/test-cache",
					"id":     "b/test-cache/roles/storage.objectCreator/user:user.name@email.com",
					"member": "user:user.name@email.com",
					"role":   "roles/storage.objectCreator",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_storage_bucket_iam_member.test",
				ResourceID:      "b/test-cache roles/storage.objectCreator user:user.name@email.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_storage_bucket_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_storage_bucket_iam_binding.test",
				Type:    "google_storage_bucket_iam_binding",
				AttributeValues: map[string]any{
					"bucket": "b/test-company-com",
					"id":     "b/test-company-com/roles/storage.legacyObjectReader",
					"role":   "roles/storage.legacyObjectReader",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_storage_bucket_iam_binding.test",
				ResourceID:      "b/test-company-com roles/storage.legacyObjectReader",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_compute_subnetwork_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_compute_subnetwork_iam_binding.test",
				Type:    "google_compute_subnetwork_iam_binding",
				AttributeValues: map[string]any{
					"id":         "projects/project/regions/us-north1/subnetworks/test/roles/compute.networkUser",
					"role":       "roles/compute.networkUser",
					"subnetwork": "projects/project/regions/us-north1/subnetworks/test",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_compute_subnetwork_iam_binding.test",
				ResourceID:      "projects/project/regions/us-north1/subnetworks/test roles/compute.networkUser",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_organization_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_organization_iam_member.test",
				Type:    "google_organization_iam_member",
				AttributeValues: map[string]any{
					"id":     "123456789/roles/viewer/serviceAccount:test@prod-eu4.iam.gserviceaccount.com",
					"member": "serviceAccount:test@prod-eu4.iam.gserviceaccount.com",
					"org_id": "123456789",
					"role":   "roles/viewer",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_organization_iam_member.test",
				ResourceID:      "123456789 roles/viewer serviceAccount:test@prod-eu4.iam.gserviceaccount.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_project_iam_custom_role",
			terraformResource: parser.TerraformResource{
				Address: "google_project_iam_custom_role.test",
				Type:    "google_project_iam_custom_role",
				AttributeValues: map[string]any{
					"id":      "projects/project/roles/osLoginProjectGet_23i4po",
					"project": "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_project_iam_custom_role.test",
				ResourceID:      "project projects/project/roles/osLoginProjectGet_23i4po",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_pubsub_topic_iam_binding",
			terraformResource: parser.TerraformResource{
				Address: "google_pubsub_topic_iam_binding.test",
				Type:    "google_pubsub_topic_iam_binding",
				AttributeValues: map[string]any{
					"id":      "projects/project/topics/subscription-test/roles/pubsub.publisher",
					"project": "project",
					"role":    "roles/pubsub.publisher",
					"topic":   "projects/project/topics/subscription-test",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_pubsub_topic_iam_binding.test",
				ResourceID:      "projects/project/topics/subscription-test roles/pubsub.publisher",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_pubsub_topic_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_pubsub_topic_iam_member.test",
				Type:    "google_pubsub_topic_iam_member",
				AttributeValues: map[string]any{
					"id":      "projects/project/topics/service-data-logs-export/roles/pubsub.publisher/serviceAccount:service@gcp-sa-logging.iam.gserviceaccount.com",
					"member":  "serviceAccount:service@gcp-sa-logging.iam.gserviceaccount.com",
					"project": "project",
					"role":    "roles/pubsub.publisher",
					"topic":   "projects/project/topics/service-data-logs-export",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_pubsub_topic_iam_member.test",
				ResourceID:      "projects/project/topics/service-data-logs-export roles/pubsub.publisher serviceAccount:service@gcp-sa-logging.iam.gserviceaccount.com",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_pubsub_subscription_iam_member",
			terraformResource: parser.TerraformResource{
				Address: "google_pubsub_subscription_iam_member.test",
				Type:    "google_pubsub_subscription_iam_member",
				AttributeValues: map[string]any{
					"id":           "projects/project/subscriptions/subscription-test/roles/pubsub.subscriber/serviceAccount:service@gcp-sa-pubsub.iam.gserviceaccount.com",
					"member":       "serviceAccount:service@gcp-sa-pubsub.iam.gserviceaccount.com",
					"project":      "project",
					"role":         "roles/pubsub.subscriber",
					"subscription": "subscription-test",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_pubsub_subscription_iam_member.test",
				ResourceID:      "projects/project/subscriptions/subscription-test/roles/pubsub.subscriber/serviceAccount:service@gcp-sa-pubsub.iam.gserviceaccount.com",
				SupportsImport:  false,
			},
		},
		{
			name: "For google_resource_manager_lien",
			terraformResource: parser.TerraformResource{
				Address: "google_resource_manager_lien.lien",
				Type:    "google_resource_manager_lien",
				AttributeValues: map[string]any{
					"name":   "p1234567890-apa12345-ab12-2727-bc34-d1234567890",
					"origin": "project-factory",
					"parent": "projects/123456890",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_resource_manager_lien.lien",
				ResourceID:      "123456890/p1234567890-apa12345-ab12-2727-bc34-d1234567890",
				SupportsImport:  true,
			},
		},
		{
			name: "For google_monitoring_uptime_check_config",
			terraformResource: parser.TerraformResource{
				Address: "google_monitoring_uptime_check_config.test",
				Type:    "google_monitoring_uptime_check_config",
				AttributeValues: map[string]any{
					"id":      "projects/project/uptimeCheckConfigs/uptimecheckconfigtest",
					"name":    "projects/project/uptimeCheckConfigs/uptimecheckconfigtest",
					"project": "project",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "google_monitoring_uptime_check_config.test",
				ResourceID:      "project projects/project/uptimeCheckConfigs/uptimecheckconfigtest",
				SupportsImport:  true,
			},
		},
		{
			name: "For everything else",
			terraformResource: parser.TerraformResource{
				Address: "example.address",
				Type:    "example_type",
				AttributeValues: map[string]any{
					"id": "test_id",
				},
			},
			expected: TerraformImport{
				ResourceAddress: "example.address",
				ResourceID:      "test_id",
				SupportsImport:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := computeTerraformImportForResource(tt.terraformResource)
			require.Equal(t, tt.expected, actual)
		})
	}
}
