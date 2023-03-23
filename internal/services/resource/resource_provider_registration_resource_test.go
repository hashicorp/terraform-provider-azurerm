package resource_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this can be moved up a level when all the others are

type ResourceProviderRegistrationResource struct{}

func TestAccResourceProviderRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("Microsoft.BlockchainTokens"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceProviderRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("Microsoft.Marketplace"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport("Microsoft.Marketplace")
		}),
	})
}

func TestAccResourceProviderRegistration_feature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				// Last error may cause resource provider still in `Registered` status.Need to unregister it before a new test.
				if err := r.unRegisterProviders("Microsoft.ApiSecurity"); err != nil {
					t.Fatalf("Failed to reset feature registration with error: %+v", err)
				}
			},
			Config: r.multiFeature(true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ResourceProviderRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	resp, err := client.Resource.ProvidersClient.Get(ctx, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("Bad: Get on ProvidersClient: %+v", err)
	}

	return utils.Bool(resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Registered")), nil
}

func (ResourceProviderRegistrationResource) basic(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
  lifecycle {
    ignore_changes = [feature]
  }
}
`, name)
}

func (r ResourceProviderRegistrationResource) requiresImport(name string) string {
	template := r.basic(name)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_provider_registration" "import" {
  name = azurerm_resource_provider_registration.test.name
}
`, template)
}

func (r ResourceProviderRegistrationResource) unRegisterProviders(resourceProviders ...string) error {
	client, err := testclient.Build()
	if err != nil {
		return fmt.Errorf("building client: %+v", err)
	}

	for _, rp := range resourceProviders {
		if err = r.unRegisterProvider(client, rp); err != nil {
			return err
		}
	}

	return nil
}

func (r ResourceProviderRegistrationResource) unRegisterProvider(client *clients.Client, resourceProvider string) error {
	ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(30*time.Minute))
	defer cancel()

	providersClient := client.Resource.ProvidersClient
	provider, err := providersClient.Get(ctx, resourceProvider, "")
	if err != nil {
		return fmt.Errorf("retrieving Resource Provider %q: %+v", resourceProvider, err)
	}

	if provider.RegistrationState == nil {
		return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was nil", resourceProvider)
	}

	if !strings.EqualFold(*provider.RegistrationState, "Registered") {
		return nil
	}

	if _, err := providersClient.Unregister(ctx, resourceProvider); err != nil {
		return fmt.Errorf("unregistering Resource Provider %q: %+v", resourceProvider, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("could not retrieve context deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Processing"},
		Target:  []string{"Unregistered"},
		Refresh: func() (interface{}, string, error) {
			resp, err := providersClient.Get(ctx, resourceProvider, "")
			if err != nil {
				return resp, "Failed", err
			}

			if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Unregistered") {
				return resp, "Unregistered", nil
			}

			return resp, "Processing", nil
		},
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Resource Provider %q to become unregistered: %+v", resourceProvider, err)
	}

	return nil
}

func (ResourceProviderRegistrationResource) multiFeature(registered1 bool, registered2 bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = "Microsoft.ApiSecurity"
  feature {
    name       = "PP2CanaryAccessDEV"
    registered = %t
  }
  feature {
    name       = "PP3CanaryAccessDEV"
    registered = %t
  }
}
`, registered1, registered2)
}
