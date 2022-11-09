package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutomanageConfigurationProfileHCIAssignmentResource struct{}

func TestAccAutomanageConfigurationProfileHCIAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_hci_assignment", "test")
	r := AutomanageConfigurationProfileHCIAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutomanageConfigurationProfileHCIAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_hci_assignment", "test")
	r := AutomanageConfigurationProfileHCIAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileHCIAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationProfileHCIAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automanage.ConfigurationProfileHCIAssignmentClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}
	return utils.Bool(true), nil
}

func (r AutomanageConfigurationProfileHCIAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_stack_hci_cluster" "test" {
  name                = "acctest-azshci-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  client_id           = data.azurerm_client_config.current.client_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AutomanageConfigurationProfileHCIAssignmentResource) requiresImport(data acceptance.TestData) string {
	config := r.complete(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile_hci_assignment" "import" {
  name                     = azurerm_automanage_configuration_profile_hci_assignment.test.name
  resource_group_name      = azurerm_automanage_configuration_profile_hci_assignment.test.resource_group_name
  cluster_name             = azurerm_automanage_configuration_profile_hci_assignment.test.cluster_name
  configuration_profile_id = azurerm_automanage_configuration_profile_hci_assignment.test.configuration_profile_id
}
`, config)
}

func (r AutomanageConfigurationProfileHCIAssignmentResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile_hci_assignment" "test" {
  name                     = "default"
  resource_group_name      = azurerm_resource_group.test.name
  cluster_name             = azurerm_stack_hci_cluster.test.name
  configuration_profile_id = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
`, template)
}
