package tfimportgen

import (
	"github.com/kishaningithub/tf-import-gen/pkg/internal/parser"
	"github.com/stretchr/testify/require"
	"testing"
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
