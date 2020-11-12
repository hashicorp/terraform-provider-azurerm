package acceptance

import (
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

var once sync.Once

func EnsureProvidersAreInitialised() {
	if !enableBinaryTesting {
		// NOTE: (@tombuildsstuff) - opting-out of Binary Testing for the moment
		os.Setenv("TF_DISABLE_BINARY_TESTING", "true")
	}

	once.Do(func() {
		azureProvider := provider.TestAzureProvider().(*schema.Provider)

		AzureProvider = azureProvider
		SupportedProviders = map[string]terraform.ResourceProvider{
			"azurerm": azureProvider,
			"azuread": azuread.Provider().(*schema.Provider),
		}

		// NOTE: (@tombuildsstuff) - intentionally not calling these as Binary Testing
		// is Disabled
		// binarytestfuntime.UseBinaryDriver("azurerm", provider.TestAzureProvider)
		// binarytestfuntime.UseBinaryDriver("azuread", azuread.Provider)
	})
}
