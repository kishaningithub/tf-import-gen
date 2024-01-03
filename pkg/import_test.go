package tfimportgen_test

import (
	"fmt"
	"github.com/kishaningithub/tf-import-gen/pkg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImport_ShouldSerializeAsTerraformImportStatements(t *testing.T) {
	tfImport := tfimportgen.TerraformImport{
		ResourceAddress: "aws_glue_catalog_database.test_db",
		ResourceID:      "id_test_db",
		SupportsImport:  true,
	}

	expectedResult := `import {
  to = aws_glue_catalog_database.test_db
  id = "id_test_db"
}
`

	require.Equal(t, expectedResult, tfImport.String())
	require.Equal(t, expectedResult, fmt.Sprint(tfImport))
	require.Equal(t, expectedResult, fmt.Sprintf("%s", tfImport))
	require.Equal(t, expectedResult, fmt.Sprintf("%v", tfImport))
}

func TestImports_ShouldSerializeAsMultipleTerraformImportStatements(t *testing.T) {
	imports := tfimportgen.TerraformImports{
		{
			ResourceAddress: "aws_glue_catalog_database.test_db",
			ResourceID:      "id_test_db",
			SupportsImport:  true,
		},
		{
			ResourceAddress: "aws_iam_instance_profile.test_instance_profile",
			ResourceID:      "id_test_instance_profile",
			SupportsImport:  true,
		},
	}

	expectedResult := `import {
  to = aws_glue_catalog_database.test_db
  id = "id_test_db"
}

import {
  to = aws_iam_instance_profile.test_instance_profile
  id = "id_test_instance_profile"
}

`

	require.Equal(t, expectedResult, imports.String())
	require.Equal(t, expectedResult, fmt.Sprint(imports))
	require.Equal(t, expectedResult, fmt.Sprintf("%s", imports))
	require.Equal(t, expectedResult, fmt.Sprintf("%v", imports))
}

func TestImport_ShouldGenerateHelpfulMessageWhenResourceDoesNotSupportImport(t *testing.T) {
	tfImport := tfimportgen.TerraformImport{
		SupportsImport:  false,
		ResourceAddress: "aws_alb_target_group_attachment.test_alb_target_group_attachment",
		ResourceID:      "id_test_alb_target_group_attachment",
	}

	expectedResult := `# resource "aws_alb_target_group_attachment.test_alb_target_group_attachment" with identifier "id_test_alb_target_group_attachment" does not support import operation. Kindly refer resource documentation for more info.`

	require.Equal(t, fmt.Sprintln(expectedResult), fmt.Sprint(tfImport))
}
