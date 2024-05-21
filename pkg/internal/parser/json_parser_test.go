package parser

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTerraformStateJsonParserAddressComputationForRootResources(t *testing.T) {
	tests := []struct {
		name                    string
		inputTerraformStateJson string
		computedAddress         string
	}{
		{
			name: "root resource without index",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "values": {
					"root_module": {
					  "resources": [
						{
						  "address": "aws_glue_catalog_database.test_db",
						  "mode": "managed",
						  "type": "aws_glue_catalog_database",
						  "name": "test_db",
						  "provider_name": "aws",
						  "values": {
							"id": "id_test_db"
						  }
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "aws_glue_catalog_database.test_db",
		},
		{
			name: "root resource with index as integer and included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "resources": [
						{
						  "address": "aws_glue_catalog_database.test_db[0]",
						  "mode": "managed",
						  "type": "aws_glue_catalog_database",
						  "name": "test_db",
                          "index": 0,
						  "provider_name": "aws",
						  "values": {
							"id": "id_test_db"
						  }
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "aws_glue_catalog_database.test_db[0]",
		},
		{
			name: "root resource with index as integer and not included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "resources": [
						{
						  "address": "aws_glue_catalog_database.test_db",
						  "mode": "managed",
						  "type": "aws_glue_catalog_database",
						  "name": "test_db",
                          "index": 0,
						  "provider_name": "aws",
						  "values": {
							"id": "id_test_db"
						  }
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "aws_glue_catalog_database.test_db[0]",
		},
		{
			name: "root resource with index as string and included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "resources": [
						{
						  "address": "aws_glue_catalog_database.test_db[\"one\"]",
						  "mode": "managed",
						  "type": "aws_glue_catalog_database",
						  "name": "test_db",
                          "index": "one",
						  "provider_name": "aws",
						  "values": {
							"id": "id_test_db"
						  }
						}
					  ]
					}
				  }
				}
`,
			computedAddress: `aws_glue_catalog_database.test_db["one"]`,
		},
		{
			name: "root resource with index as string and not included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "resources": [
						{
						  "address": "aws_glue_catalog_database.test_db",
						  "mode": "managed",
						  "type": "aws_glue_catalog_database",
						  "name": "test_db",
                          "index": "one",
						  "provider_name": "aws",
						  "values": {
							"id": "id_test_db"
						  }
						}
					  ]
					}
				  }
				}
`,
			computedAddress: `aws_glue_catalog_database.test_db["one"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewTerraformStateJsonParser(bytes.NewBufferString(tt.inputTerraformStateJson))
			actualResources, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, tt.computedAddress, actualResources[0].Address)
		})
	}
}

func TestTerraformStateJsonParserDoesNotAddModuleNameTwice(t *testing.T) {
	tests := []struct {
		name                    string
		inputTerraformStateJson string
		computedAddress         string
	}{
		{
			name: "1 level repeated module address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "module.test_mwaa.aws_iam_policy.test_mwaa_permissions",
							  "mode": "managed",
							  "type": "aws_iam_policy",
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewTerraformStateJsonParser(bytes.NewBufferString(tt.inputTerraformStateJson))
			actualResources, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, tt.computedAddress, actualResources[0].Address)
		})
	}
}

func TestTerraformStateJsonParserAddressComputationForNestedResources(t *testing.T) {
	tests := []struct {
		name                    string
		inputTerraformStateJson string
		computedAddress         string
	}{
		{
			name: "nested resource without index",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "aws_iam_policy.test_mwaa_permissions",
							  "mode": "managed",
							  "type": "aws_iam_policy",
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions",
		},
		{
			name: "nested resource with index as integer and included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "aws_iam_policy.test_mwaa_permissions[0]",
							  "mode": "managed",
							  "type": "aws_iam_policy",
                              "index": 0,
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions[0]",
		},
		{
			name: "nested resource with index as integer and not included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "aws_iam_policy.test_mwaa_permissions",
							  "mode": "managed",
							  "type": "aws_iam_policy",
                              "index": 0,
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions[0]",
		},
		{
			name: "nested resource with index as string and included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "aws_iam_policy.test_mwaa_permissions[\"one\"]",
							  "mode": "managed",
							  "type": "aws_iam_policy",
                              "index": "one",
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: `module.test_mwaa.aws_iam_policy.test_mwaa_permissions["one"]`,
		},
		{
			name: "nested resource with index as string and not included in address",
			inputTerraformStateJson: `
				{
				  "format_version": "0.1",
				  "terraform_version": "0.12.31",
				  "values": {
					"root_module": {
					  "child_modules": [
						{
						  "resources": [
							{
							  "address": "aws_iam_policy.test_mwaa_permissions",
							  "mode": "managed",
							  "type": "aws_iam_policy",
                              "index": "one",
							  "name": "test_mwaa_permissions",
							  "provider_name": "aws",
							  "values": {
								"id": "id_test_mwaa_permissions"
							  }
							}
						  ],
						  "address": "module.test_mwaa"
						}
					  ]
					}
				  }
				}
`,
			computedAddress: `module.test_mwaa.aws_iam_policy.test_mwaa_permissions["one"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewTerraformStateJsonParser(bytes.NewBufferString(tt.inputTerraformStateJson))
			actualResources, err := parser.Parse()
			require.NoError(t, err)
			require.Equal(t, tt.computedAddress, actualResources[0].Address)
		})
	}
}
