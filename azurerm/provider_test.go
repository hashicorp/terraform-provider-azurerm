package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	azureProvider := Provider().(*schema.Provider)

	testAccProvider = azureProvider
	acceptance.AzureProvider = azureProvider

	// NOTE: these /cannot/ be simplified into a single shared variable (tried, it causes a nil-slice)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
	acceptance.SupportedProviders = map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}
