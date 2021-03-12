package acceptance

import (
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/acctest"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

var once sync.Once

func EnsureProvidersAreInitialised() {
	// require reattach testing is enabled
	os.Setenv("TF_ACCTEST_REATTACH", "1")

	once.Do(func() {
		acctest.UseBinaryDriver("azurerm", provider.TestAzureProvider)
		acctest.UseBinaryDriver("azuread", azuread.Provider)
	})
}
