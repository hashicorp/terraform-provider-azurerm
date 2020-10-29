package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// NOTE: when Binary Testing is enabled the Check functions will need to build a client rather than relying on the
// shared one. For the moment we init the shared client when Binary Testing is Enabled & Disabled - but this needs
// fixing when we move to Binary Testing so that we can test across provider instances
var enableBinaryTesting = false

func (td TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		Steps:    steps,
	}

	td.runAcceptanceTest(t, testCase)
}

func (td TestData) ResourceTest(t *testing.T, testResource types.TestResource, steps []resource.TestStep) {
	testCase := resource.TestCase{
		PreCheck: func() { PreCheck(t) },
		CheckDestroy: func(s *terraform.State) error {
			client := buildClient()
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.resourceLabel)(s)
		},
		Steps: steps,
	}

	td.runAcceptanceTest(t, testCase)
}

func buildClient() *clients.Client {
	if enableBinaryTesting {
		// TODO: we can likely cache the built up auth object and return a new client instance
	}

	return AzureProvider.Meta().(*clients.Client)
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	if enableBinaryTesting {
		testCase.ProviderFactories = map[string]terraform.ResourceProviderFactory{
			// TODO: switch this out for dynamic initialization
			"azurerm": terraform.ResourceProviderFactoryFixed(AzureProvider),
		}
		//} else {
	}
	testCase.Providers = SupportedProviders

	resource.ParallelTest(t, testCase)
}
