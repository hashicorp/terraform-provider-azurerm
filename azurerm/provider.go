package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func Provider() terraform.ResourceProvider {
	return provider.AzureProvider()
}
