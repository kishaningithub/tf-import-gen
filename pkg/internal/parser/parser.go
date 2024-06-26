package parser

import (
	"strings"
)

type TerraformResource struct {
	Address         string
	Type            string
	Index           any
	AttributeValues map[string]any
}

type TerraformResources []TerraformResource

func (resources TerraformResources) FilterByAddresses(addresses []string) TerraformResources {
	var filteredResources TerraformResources
	for _, resource := range resources {
		for _, address := range addresses {
			if strings.HasPrefix(resource.Address, address) {
				filteredResources = append(filteredResources, resource)
				break
			}
		}
	}
	return filteredResources
}

type TerraformStateParser interface {
	Parse() (TerraformResources, error)
}
