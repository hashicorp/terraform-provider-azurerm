package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

type SiteRecoveryProtectionContainerMappingResource struct {
}

func TestAccSiteRecoveryProtectionContainerMapping_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_protection_container_mapping", "test")
	r := SiteRecoveryProtectionContainerMappingResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryProtectionContainerMappingResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-recovery-%d-1"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_fabric" "test1" {
  resource_group_name = azurerm_resource_group.test1.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric1-%d"
  location            = azurerm_resource_group.test1.location
}

resource "azurerm_site_recovery_fabric" "test2" {
  resource_group_name = azurerm_resource_group.test1.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric2-%d"
  location            = "%s"
  depends_on          = [azurerm_site_recovery_fabric.test1]
}

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test1.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test1.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  name                 = "acctest-protection-cont2-%d"
}

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = azurerm_resource_group.test1.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.test.name
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test1.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t SiteRecoveryProtectionContainerMappingResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectionContainerMappings"]

	resp, err := clients.RecoveryServices.ContainerMappingClient(resGroup, vaultName).Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery protection container mapping (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
