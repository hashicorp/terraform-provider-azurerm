package providerjson

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func (p *ProviderJSON) DataSources() []terraform.DataSource {
	s := schema.Provider(*p)
	return s.DataSources()
}

func (p *ProviderJSON) Resources() []terraform.ResourceType {
	s := schema.Provider(*p)
	return s.Resources()
}
