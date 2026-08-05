// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SiteRecoveryProtectionContainerMappingResource struct{}

func TestAccSiteRecoveryProtectionContainerMapping_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_protection_container_mapping", "test")
	r := SiteRecoveryProtectionContainerMappingResource{}

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

func TestAccSiteRecoveryProtectionContainerMapping_withSystemAssignedAutoUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_protection_container_mapping", "test")
	r := SiteRecoveryProtectionContainerMappingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoUpdateExtension(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoUpdateExtension(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoUpdateExtension(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoUpdateExtensionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SiteRecoveryProtectionContainerMappingResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[2]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[1]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_fabric" "test1" {
  resource_group_name = azurerm_resource_group.test1.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric1-%[1]d"
  location            = azurerm_resource_group.test1.location
}

resource "azurerm_site_recovery_fabric" "test2" {
  resource_group_name = azurerm_resource_group.test1.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric2-%[1]d"
  location            = "%[3]s"
  depends_on          = [azurerm_site_recovery_fabric.test1]
}

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test1.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%[1]d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test1.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  name                 = "acctest-protection-cont2-%[1]d"
}

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = azurerm_resource_group.test1.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.test.name
  name                                                 = "acctest-policy-%[1]d"
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SiteRecoveryProtectionContainerMappingResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test1.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%[2]d"
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryProtectionContainerMappingResource) autoUpdateExtension(data acceptance.TestData, enabled bool) string {
	automaticUpdate := ""
	if enabled {
		automaticUpdate = `
  automatic_update {
    automation_account_id = azurerm_automation_account.test.id
    authentication_type   = "SystemAssignedIdentity"
  }
`
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%[2]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name

  sku_name = "Basic"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test1.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%[2]d"

%[3]s
}
`, r.template(data), data.RandomInteger, automaticUpdate)
}

func (r SiteRecoveryProtectionContainerMappingResource) autoUpdateExtensionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%[2]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name

  sku_name = "Basic"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_automation_account" "test2" {
  name                = "acctestAutomation-2-%[2]d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name

  sku_name = "Basic"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test1.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%[2]d"
  automatic_update {
    automation_account_id = azurerm_automation_account.test2.id
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryProtectionContainerMappingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ContainerMappingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery protection container mapping (%s): %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading site recovery protection container mapping (%s): model is nil", id.String())
	}

	return pointer.To(model.Id != nil), nil
}
