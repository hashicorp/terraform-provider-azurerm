package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SynapsePrivateLinkHubResource struct{}

func TestAccSynapsePrivateLinkHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_private_link_hub", "test")
	r := SynapsePrivateLinkHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapsePrivateLinkHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_private_link_hub", "test")
	r := SynapsePrivateLinkHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapsePrivateLinkHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_private_link_hub", "test")
	r := SynapsePrivateLinkHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withUpdateFields(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test2"),
			),
		},
	})
}

func (r SynapsePrivateLinkHubResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PrivateLinkHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.PrivateLinkHubsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SynapsePrivateLinkHubResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_private_link_hub" "test" {
  name                = "acctestsw%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r SynapsePrivateLinkHubResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_private_link_hub" "import" {
  name                = azurerm_synapse_private_link_hub.test.name
  resource_group_name = azurerm_synapse_private_link_hub.test.resource_group_name
  location            = azurerm_synapse_private_link_hub.test.location
}
`, config)
}

func (r SynapsePrivateLinkHubResource) withUpdateFields(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_private_link_hub" "test" {
  name                = "acctestsw%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger)
}

func (r SynapsePrivateLinkHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
