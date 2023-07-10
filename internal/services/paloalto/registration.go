package paloalto

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) Name() string {
	return "Palo Alto"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LocalRuleStack{},
		LocalRuleStackCertificate{},
		LocalRuleStackRule{},
		NextGenerationFirewall{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Palo Alto",
	}
}
