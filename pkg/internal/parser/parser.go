package parser

import (
	"strings"
)

type TerraformResource struct {
	Address         string
	Type            string
	AttributeValues map[string]any
}

type TerraformResources []TerraformResource

func (resources TerraformResources) FilterByAddresses(addresses []string) TerraformResources {
	var filteredResources TerraformResources
	for _, resource := range resources {
		for _, address := range addresses {
			if strings.HasPrefix(resource.Address, address) {
				if !filteredResources.contains(resource) {
					filteredResources = append(filteredResources, resource)
				}
			}
		}
	}
	return filteredResources
}

func (resources TerraformResources) contains(r TerraformResource) bool {
	for _, resource := range resources {
		if resource.Address == r.Address {
			return true
		}
	}
	return false
}

type TerraformStateParser interface {
	Parse() (TerraformResources, error)
}
