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
	supportedProviders := map[string]terraform.ResourceProvider{
		"azurerm": testAccProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}

	// TODO: these can be de-duped once this is relocated
	testAccProvider = azureProvider
	acceptance.AzureProvider = azureProvider
	testAccProviders = supportedProviders
	acceptance.SupportedProviders = supportedProviders
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}
