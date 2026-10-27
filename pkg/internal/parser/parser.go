package parser

import (
	"strings"

	"github.com/samber/lo"
)

type TerraformResource struct {
	Address         string
	Type            string
	Index           any
	AttributeValues map[string]any
}

type TerraformResources []TerraformResource

func (resources TerraformResources) FilterByAddresses(addresses []string) TerraformResources {
	return lo.Filter(resources, func(resource TerraformResource, _ int) bool {
		return lo.ContainsBy(addresses, func(address string) bool {
			return strings.HasPrefix(resource.Address, address)
		})
	})
}

type TerraformStateParser interface {
	Parse() (TerraformResources, error)
}
