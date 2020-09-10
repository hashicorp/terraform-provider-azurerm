package acceptance

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

var once sync.Once

func EnsureProvidersAreInitialised() {
	// NOTE: (@tombuildsstuff) - opting-out of Binary Testing for the moment
	// os.Setenv("TF_DISABLE_BINARY_TESTING", "true")

	//once.Do(func() {
		// azureProvider := provider.TestAzureProvider().(*schema.Provider)
		// azureProvider2 := provider.TestAzureProvider().(*schema.Provider)

		// AzureProvider = azureProvider
		SupportedProviders = map[string]terraform.ResourceProviderFactory{
			"azurerm":  terraform.ResourceProviderFactoryFixed(provider.TestAzureProvider()),
			"azurerm2": terraform.ResourceProviderFactoryFixed(provider.TestAzureProvider()),
			// "azuread":  azuread.Provider().(*schema.Provider),
		}

		// NOTE: (@tombuildsstuff) - intentionally not calling these as Binary Testing
		// is Disabled
		// binarytestfuntime.UseBinaryDriver("azurerm", provider.TestAzureProvider)
		//binarytestfuntime.UseBinaryDriver("azuread", azuread.Provider)
	//})
}
