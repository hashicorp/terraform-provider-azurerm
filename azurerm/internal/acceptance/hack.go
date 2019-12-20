package acceptance

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
)

func CustomInit(provider *schema.Provider) {
	AzureProvider = provider

	SupportedProviders = map[string]terraform.ResourceProvider{
		"azurerm": provider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
}
