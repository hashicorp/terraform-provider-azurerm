package acceptance

import (
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/testclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

var once sync.Once

func EnsureProvidersAreInitialised() {
	if !testclient.EnableBinaryTesting {
		os.Setenv("TF_DISABLE_BINARY_TESTING", "true")
	} else {
		// require reattach testing is enabled
		os.Setenv("TF_ACCTEST_REATTACH", "1")
	}

	once.Do(func() {
		if !testclient.EnableBinaryTesting {
			azureProvider := provider.TestAzureProvider().(*schema.Provider)
			testclient.AzureProvider = azureProvider
			testclient.SupportedProviders = map[string]terraform.ResourceProvider{
				"azurerm": azureProvider,
				"azuread": azuread.Provider().(*schema.Provider),
			}
		} else {
			acctest.UseBinaryDriver("azurerm", provider.TestAzureProvider)
			acctest.UseBinaryDriver("azuread", azuread.Provider)
		}
	})
}
