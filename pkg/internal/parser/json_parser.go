package parser

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"io"
)

var (
	_ TerraformStateParser = TerraformStateJsonParser{}
)

type TerraformStateJsonParser struct {
	reader io.Reader
}

func NewTerraformStateJsonParser(reader io.Reader) TerraformStateParser {
	return TerraformStateJsonParser{
		reader: reader,
	}
}

func (parser TerraformStateJsonParser) Parse() (TerraformResources, error) {
	tfStateJsonBytes, err := io.ReadAll(parser.reader)
	if err != nil {
		return nil, err
	}
	var state tfjson.State
	err = json.Unmarshal(tfStateJsonBytes, &state)
	if err != nil {
		return nil, err
	}
	return parser.parseResourcesWithinModule(state.Values.RootModule), nil
}

func (parser TerraformStateJsonParser) parseResourcesWithinModule(module *tfjson.StateModule) []TerraformResource {
	var allResources []TerraformResource
	allResources = append(allResources, parser.parseResources(module.Resources, module.Address)...)
	for _, childModule := range module.ChildModules {
		allResources = append(allResources, parser.parseResourcesWithinModule(childModule)...)
	}
	return allResources
}

func (parser TerraformStateJsonParser) parseResources(resources []*tfjson.StateResource, moduleAddress string) []TerraformResource {
	var resourceImportModel []TerraformResource
	for _, resource := range resources {
		if resource.Mode != tfjson.ManagedResourceMode {
			continue
		}
		resourceImportModel = append(resourceImportModel, TerraformResource{
			Address:         parser.computeResourceAddress(moduleAddress, resource.Address),
			Type:            resource.Type,
			AttributeValues: resource.AttributeValues,
		})
	}
	return resourceImportModel
}

func (parser TerraformStateJsonParser) computeResourceAddress(moduleAddress string, resourceAddress string) string {
	if len(moduleAddress) == 0 {
		return resourceAddress
	}
	return fmt.Sprintf("%s.%s", moduleAddress, resourceAddress)
}
