package internal_test

import (
	"fmt"
	tfimportgen "github.com/kishaningithub/tf-import-gen/pkg/internal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImport_ShouldSerializeAsTerraformImportStatements(t *testing.T) {
	tfImport := tfimportgen.TerraformImport{
		ResourceAddress: "aws_glue_catalog_database.test_db",
		ResourceID:      "id_test_db",
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
		},
		{
			ResourceAddress: "aws_iam_instance_profile.test_instance_profile",
			ResourceID:      "id_test_instance_profile",
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
