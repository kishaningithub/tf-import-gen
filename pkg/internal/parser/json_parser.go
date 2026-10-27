package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/samber/lo"
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
	resourcesWithinModule := parser.parseResources(module.Resources, module.Address)
	resourcesWithinChildModules := lo.FlatMap(module.ChildModules, func(childModule *tfjson.StateModule, _ int) []TerraformResource {
		return parser.parseResourcesWithinModule(childModule)
	})
	return append(resourcesWithinModule, resourcesWithinChildModules...)
}

func (parser TerraformStateJsonParser) parseResources(resources []*tfjson.StateResource, moduleAddress string) []TerraformResource {
	managedResources := lo.Filter(resources, func(resource *tfjson.StateResource, _ int) bool {
		return resource.Mode == tfjson.ManagedResourceMode
	})
	return lo.Map(managedResources, func(resource *tfjson.StateResource, _ int) TerraformResource {
		return TerraformResource{
			Address:         parser.computeResourceAddressIncludingModule(moduleAddress, resource),
			Type:            resource.Type,
			AttributeValues: resource.AttributeValues,
		}
	})
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
