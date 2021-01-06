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

// lintignore:AT001
func (td TestData) DataSourceTest(t *testing.T, steps []resource.TestStep) {
	// DataSources don't need a check destroy - however since this is a wrapper function
	// and not matching the ignore pattern `XXX_data_source_test.go`, this needs to be explicitly opted out
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
			return helpers.CheckDestroyedFunc(client, testResource, td.ResourceType, td.ResourceName)(s)
		},
		Steps: steps,
	}

	td.runAcceptanceTest(t, testCase)
}

func RunTestsInSequence(t *testing.T, tests map[string]map[string]func(t *testing.T)) {
	for group, m := range tests {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func buildClient() *clients.Client {
	// if enableBinaryTesting {
	//   TODO: build up a client on demand
	//   NOTE: this'll want caching/a singleton, and likely RP registration etc disabled, since otherwise this'll become
	//   		 extremely expensive - and this doesn't need access to the provider feature toggles
	// }

	return AzureProvider.Meta().(*clients.Client)
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	if enableBinaryTesting {
		testCase.ProviderFactories = map[string]terraform.ResourceProviderFactory{
			// TODO: switch this out for dynamic initialization?
			"azurerm": terraform.ResourceProviderFactoryFixed(AzureProvider),
		}
	}
	testCase.Providers = SupportedProviders

	resource.ParallelTest(t, testCase)
}
