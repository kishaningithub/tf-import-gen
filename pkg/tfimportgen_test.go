package tfimportgen_test

import (
	tfimportgen "github.com/kishaningithub/tf-import-gen/pkg"
	"github.com/kishaningithub/tf-import-gen/pkg/internal"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_GenerateImports_ShouldGenerateImportsForAllResourcesWhenNoFiltersAreGiven(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected internal.TerraformImports
	}{
		{
			name:     "only root resources",
			filePath: "testdata/only_root_resources.json",
			expected: internal.TerraformImports{
				{
					ResourceAddress: "aws_glue_catalog_database.test_db",
					ResourceID:      "id_test_db",
				},
				{
					ResourceAddress: "aws_iam_instance_profile.test_instance_profile",
					ResourceID:      "id_test_instance_profile",
				},
			},
		},
		{
			name:     "resources in child module",
			filePath: "testdata/resources_in_child_module.json",
			expected: internal.TerraformImports{
				{
					ResourceAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions",
					ResourceID:      "id_test_mwaa_permissions",
				},
				{
					ResourceAddress: "module.test_mwaa.aws_mwaa_environment.test_airflow_env",
					ResourceID:      "id_test_airflow_env",
				},
			},
		},
		{
			name:     "resources in deeply nested child module",
			filePath: "testdata/resources_in_deeply_nested_child_module.json",
			expected: internal.TerraformImports{
				{
					ResourceAddress: "module.test_mwaa.nested1.nested2.aws_iam_policy.test_mwaa_permissions",
					ResourceID:      "id_test_mwaa_permissions",
				},
				{
					ResourceAddress: "module.test_mwaa.nested1.nested2.aws_mwaa_environment.test_airflow_env",
					ResourceID:      "id_test_airflow_env",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateJsonFile, err := os.Open(filepath.FromSlash(tt.filePath))
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = stateJsonFile.Close()
			})

			actual, err := tfimportgen.GenerateImports(stateJsonFile, "")

			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func Test_GenerateImports_ShouldGenerateImportsForResourcesForGivenAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected internal.TerraformImports
	}{
		{
			name:    "filtering by module",
			address: "module.test_mwaa",
			expected: internal.TerraformImports{
				{
					ResourceAddress: "module.test_mwaa.aws_iam_policy.test_mwaa_permissions",
					ResourceID:      "id_test_mwaa_permissions",
				},
				{
					ResourceAddress: "module.test_mwaa.aws_mwaa_environment.test_airflow_env",
					ResourceID:      "id_test_airflow_env",
				},
			},
		},
		{
			name:    "filtering by resource",
			address: "aws_glue_catalog_database.test_db",
			expected: internal.TerraformImports{
				{
					ResourceAddress: "aws_glue_catalog_database.test_db",
					ResourceID:      "id_test_db",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateJsonFile, err := os.Open(filepath.FromSlash("testdata/resources_in_root_and_child_modules.json"))
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = stateJsonFile.Close()
			})

			actual, err := tfimportgen.GenerateImports(stateJsonFile, tt.address)

			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}