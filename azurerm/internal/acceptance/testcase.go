package acceptance

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
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
			client, err := buildClient()
			if err != nil {
				return fmt.Errorf("building client: %+v", err)
			}
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

var _client *clients.Client
var clientLock = &sync.Mutex{}

func buildClient() (*clients.Client, error) {
	if enableBinaryTesting {
		clientLock.Lock()
		defer clientLock.Unlock()

		if _client == nil {
			environment, exists := os.LookupEnv("ARM_ENVIRONMENT")
			if !exists {
				environment = "public"
			}

			builder := authentication.Builder{
				SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
				ClientID:       os.Getenv("ARM_CLIENT_ID"),
				TenantID:       os.Getenv("ARM_TENANT_ID"),
				ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
				Environment:    environment,
				MetadataHost:   os.Getenv("ARM_METADATA_HOST"),

				// we intentionally only support Client Secret auth for tests (since those variables are used all over)
				SupportsClientSecretAuth: true,
			}
			config, err := builder.Build()
			if err != nil {
				return nil, fmt.Errorf("Error building ARM Client: %+v", err)
			}

			clientBuilder := clients.ClientBuilder{
				AuthConfig:               config,
				SkipProviderRegistration: true,
				TerraformVersion:         os.Getenv("TERRAFORM_CORE_VERSION"),
				Features:                 features.Default(),
				StorageUseAzureAD:        false,
			}
			client, err := clients.Build(context.TODO(), clientBuilder)
			if err != nil {
				return nil, err
			}
			_client = client
		}

		return _client, nil
	}

	return AzureProvider.Meta().(*clients.Client), nil
}

func (td TestData) runAcceptanceTest(t *testing.T, testCase resource.TestCase) {
	if enableBinaryTesting {
		testCase.DisableBinaryDriver = false //nolint:SA1019
		testCase.ProviderFactories = map[string]terraform.ResourceProviderFactory{
			"azurerm": func() (terraform.ResourceProvider, error) {
				azurerm := provider.TestAzureProvider()
				return azurerm, nil
			},
		}
	} else {
		testCase.Providers = SupportedProviders
	}

	resource.ParallelTest(t, testCase)
}
