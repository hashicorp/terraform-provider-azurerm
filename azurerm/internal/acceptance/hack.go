package acceptance

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func init() {
	azureProvider := provider.AzureProvider().(*schema.Provider)

	AzureProvider = azureProvider
	SupportedProviders = map[string]terraform.ResourceProvider{
		"azurerm": azureProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
}
