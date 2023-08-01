package automanage_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type AutoManageConfigurationHCIAssignmentResource struct{}

func TestAccAutoManageProfileHCIAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_hci_assignment", "test")
	r := AutoManageConfigurationHCIAssignmentResource{}
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

func (r AutoManageConfigurationHCIAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationHCIAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	client := clients.Automanage.HCIAssignmentClient

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Response.Response != nil), nil
}

func (r AutoManageConfigurationHCIAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-StackHCICluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_automanage_configuration_hci_assignment" "test" {
  name                = "default"
  resource_group_name = azurerm_resource_group.test.name
  configuration_id    = azurerm_automanage_configuration.test.id
  cluster_name        = azurerm_stack_hci_cluster.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
