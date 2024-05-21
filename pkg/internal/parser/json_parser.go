package parser

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"io"
	"strings"
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
			Address:         parser.computeResourceAddressIncludingModule(moduleAddress, resource),
			Type:            resource.Type,
			AttributeValues: resource.AttributeValues,
		})
	}
	return resourceImportModel
}

func (parser TerraformStateJsonParser) computeResourceAddressIncludingModule(moduleAddress string, resource *tfjson.StateResource) string {
	resourceAddress := parser.computeResourceAddress(resource)
	if len(moduleAddress) == 0 {
		return resourceAddress
	}
	if strings.HasPrefix(resourceAddress, moduleAddress) {
		return resourceAddress
	}
	return fmt.Sprintf("%s.%s", moduleAddress, resourceAddress)
}

func (parser TerraformStateJsonParser) computeResourceAddress(resource *tfjson.StateResource) string {
	if resource.Index != nil && !strings.HasSuffix(resource.Address, "]") {
		switch resource.Index.(type) {
		case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%s[%v]", resource.Address, resource.Index)
		default:
			return fmt.Sprintf("%s[%q]", resource.Address, resource.Index)
		}
	}
	return resource.Address
}
