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
