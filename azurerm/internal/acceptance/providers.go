package acceptance

import (
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

var once sync.Once

func EnsureProvidersAreInitialised() {
	// NOTE: (@tombuildsstuff) - opting-out of Binary Testing for the moment
	os.Setenv("TF_DISABLE_BINARY_TESTING", "true")

	once.Do(func() {
		azureProvider := provider.TestAzureProvider()

		AzureProvider = azureProvider
		SupportedProviders = map[string]*schema.Provider{
			"azurerm": azureProvider,

			// KEM NOTE: This provider import should be removed.
			// TODO find where it is used and make sure it works
			// under binary testing
			// "azuread": azuread.Provider(),
		}

		// NOTE: (@tombuildsstuff) - intentionally not calling these as Binary Testing
		// is Disabled
		//binarytestfuntime.UseBinaryDriver("azurerm", provider.TestAzureProvider)
		//binarytestfuntime.UseBinaryDriver("azuread", azuread.Provider)
	})
}
