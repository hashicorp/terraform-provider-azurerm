package resource_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this can be moved up a level when all the others are

type ResourceProviderRegistrationResource struct {
}

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
			Config: r.basic("Wandisco.Fusion"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport("Wandisco.Fusion")
		}),
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
