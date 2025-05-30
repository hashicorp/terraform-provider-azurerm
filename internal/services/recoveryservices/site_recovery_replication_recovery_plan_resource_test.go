// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicationRecoveryPlan struct{}

func TestAccSiteRecoveryReplicationRecoveryPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

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

func TestAccSiteRecoveryReplicationRecoveryPlan_withPreActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPreActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_withPostActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPostActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_withMultiActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultiActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// to check the actions are in the correct order
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.0.name").HasValue("testPreAction1"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.1.name").HasValue("testPreAction2"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.0.name").HasValue("testPostAction1"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.1.name").HasValue("testPostAction2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_updateWithmultiActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultiActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// to check the actions are in the correct order
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.0.name").HasValue("testPreAction1"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.1.name").HasValue("testPreAction2"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.0.name").HasValue("testPostAction1"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.1.name").HasValue("testPostAction2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateWithMultiActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// to check the actions are in the correct order
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.0.name").HasValue("testPreAction1-new"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.pre_action.1.name").HasValue("testPreAction2-new"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.0.name").HasValue("testPostAction1-new"),
				check.That(data.ResourceName).Key("boot_recovery_group.0.post_action.1.name").HasValue("testPostAction2-new"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_withEdgeZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withEdgeZones(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_withMultiBootGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withMultiBootGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("boot_recovery_group.0.replicated_protected_items.#").HasValue("1"),
				check.That(data.ResourceName).Key("boot_recovery_group.1.pre_action.0.name").HasValue("testPreAction"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSiteRecoveryReplicationRecoveryPlan_wrongActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_replication_recovery_plan", "test")
	r := SiteRecoveryReplicationRecoveryPlan{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.wrongActions(data),
			ExpectError: regexp.MustCompile("`fabric_location` must not be specified for `recovery_group` with `ManualActionDetails` type"),
		},
	})
}

func (SiteRecoveryReplicationRecoveryPlan) template(data acceptance.TestData) string {
	tags := ""
	if strings.HasPrefix(strings.ToLower(data.Client().SubscriptionID), "85b3dbca") {
		tags = `
  tags = {
    "azsecpack"                                                                = "nonprod"
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true"
  }
`
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d-1"
  location = "%[3]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-recovery-%[1]d-2"
  location = "%[4]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_fabric" "test1" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric1-%[1]d"
  location            = azurerm_resource_group.test.location
}

resource "azurerm_site_recovery_fabric" "test2" {
  resource_group_name = azurerm_resource_group.test2.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  name                = "acctest-fabric2-%[1]d"
  location            = azurerm_resource_group.test2.location
  depends_on          = [azurerm_site_recovery_fabric.test1]
}

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  name                 = "acctest-protection-cont1-%[1]d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = azurerm_resource_group.test2.name
  recovery_vault_name  = azurerm_recovery_services_vault.test.name
  recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  name                 = "acctest-protection-cont2-%[1]d"
}

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = azurerm_resource_group.test2.name
  recovery_vault_name                                  = azurerm_recovery_services_vault.test.name
  name                                                 = "acctest-policy-%[1]d"
  recovery_point_retention_in_minutes                  = 24 * 60
  application_consistent_snapshot_frequency_in_minutes = 4 * 60
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  recovery_fabric_name                      = azurerm_site_recovery_fabric.test1.name
  recovery_source_protection_container_name = azurerm_site_recovery_protection_container.test1.name
  recovery_target_protection_container_id   = azurerm_site_recovery_protection_container.test2.id
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  name                                      = "mapping-%[1]d"
}

resource "azurerm_virtual_network" "test1" {
  name                = "net-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_site_recovery_fabric.test1.location
}

resource "azurerm_subnet" "test1" {
  name                 = "snet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test1.name
  address_prefixes     = ["192.168.1.0/24"]
}

resource "azurerm_virtual_network" "test2" {
  name                = "net-%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_site_recovery_fabric.test2.location
}

resource "azurerm_subnet" "test2_1" {
  name                 = "acctest-snet-%[1]d_1"
  resource_group_name  = "${azurerm_resource_group.test2.name}"
  virtual_network_name = "${azurerm_virtual_network.test2.name}"
  address_prefixes     = ["192.168.2.0/27"]
}

resource "azurerm_subnet" "test2_2" {
  name                 = "snet-%[1]d_2"
  resource_group_name  = "${azurerm_resource_group.test2.name}"
  virtual_network_name = "${azurerm_virtual_network.test2.name}"
  address_prefixes     = ["192.168.2.32/27"]
}

resource "azurerm_site_recovery_network_mapping" "test" {
  resource_group_name         = azurerm_resource_group.test2.name
  recovery_vault_name         = azurerm_recovery_services_vault.test.name
  name                        = "mapping-%[1]d"
  source_recovery_fabric_name = azurerm_site_recovery_fabric.test1.name
  target_recovery_fabric_name = azurerm_site_recovery_fabric.test2.name
  source_network_id           = azurerm_virtual_network.test1.id
  target_network_id           = azurerm_virtual_network.test2.id
}

resource "azurerm_network_interface" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "vm-%[1]d"
    subnet_id                     = azurerm_subnet.test1.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                = "vm-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  vm_size = "Standard_B1s"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.5"
    version   = "latest"
  }

  storage_os_disk {
    name              = "disk-%[1]d"
    os_type           = "Linux"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    admin_username = "testadmin"
    admin_password = "Password1234!"
    computer_name  = "vm-%[1]d"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
  network_interface_ids = [azurerm_network_interface.test.id]

 %[5]s
}

resource "azurerm_storage_account" "test" {
  name                     = "acct%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}


resource "azurerm_site_recovery_replicated_vm" "test" {
  name                                      = "repl-%[2]d"
  resource_group_name                       = azurerm_resource_group.test2.name
  recovery_vault_name                       = azurerm_recovery_services_vault.test.name
  source_vm_id                              = azurerm_virtual_machine.test.id
  source_recovery_fabric_name               = azurerm_site_recovery_fabric.test1.name
  recovery_replication_policy_id            = azurerm_site_recovery_replication_policy.test.id
  source_recovery_protection_container_name = azurerm_site_recovery_protection_container.test1.name

  target_resource_group_id                = azurerm_resource_group.test2.id
  target_recovery_fabric_id               = azurerm_site_recovery_fabric.test2.id
  target_recovery_protection_container_id = azurerm_site_recovery_protection_container.test2.id

  managed_disk {
    disk_id                    = azurerm_virtual_machine.test.storage_os_disk[0].managed_disk_id
    staging_storage_account_id = azurerm_storage_account.test.id
    target_resource_group_id   = azurerm_resource_group.test2.id
    target_disk_type           = "Premium_LRS"
    target_replica_disk_type   = "Premium_LRS"
  }

  network_interface {
    source_network_interface_id = azurerm_network_interface.test.id
    target_subnet_name          = "snet-%[2]d_2"
  }

  depends_on = [
    azurerm_site_recovery_protection_container_mapping.test,
    azurerm_site_recovery_network_mapping.test,
  ]
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, tags)
}

func (r SiteRecoveryReplicationRecoveryPlan) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withPreActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
    pre_action {
      name                      = "testPreAction"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withPostActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
    post_action {
      name                      = "testPreAction"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withMultiActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
    pre_action {
      name                      = "testPreAction1"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    pre_action {
      name                      = "testPreAction2"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    post_action {
      name                      = "testPostAction1"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    post_action {
      name                      = "testPostAction2"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) updateWithMultiActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
    pre_action {
      name                      = "testPreAction1-new"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    pre_action {
      name                      = "testPreAction2-new"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    post_action {
      name                      = "testPostAction1-new"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }

    post_action {
      name                      = "testPostAction2-new"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
  }

  azure_to_azure_settings {
    primary_zone  = "1"
    recovery_zone = "2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withEdgeZones(data acceptance.TestData) string {
	// WestUS has an edge zone available - so hard-code to that
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
%s

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
  }

  azure_to_azure_settings {
    primary_edge_zone  = data.azurerm_extended_locations.test.extended_locations[0]
    recovery_edge_zone = data.azurerm_extended_locations.test.extended_locations[0]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) withMultiBootGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
  }

  boot_recovery_group {
    pre_action {
      name                      = "testPreAction"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) wrongActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_site_recovery_replication_recovery_plan" "test" {
  name                      = "acctest-%[2]d"
  recovery_vault_id         = azurerm_recovery_services_vault.test.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.test1.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.test2.id

  shutdown_recovery_group {}

  failover_recovery_group {}

  boot_recovery_group {
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]

    post_action {
      name                      = "testPreAction"
      type                      = "ManualActionDetails"
      fail_over_directions      = ["PrimaryToRecovery"]
      fail_over_types           = ["TestFailover"]
      manual_action_instruction = "test instruction"
      fabric_location           = "Primary"
    }
  }

}
`, r.template(data), data.RandomInteger)
}

func (r SiteRecoveryReplicationRecoveryPlan) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationRecoveryPlansClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replication plan (%s): %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return nil, fmt.Errorf("reading site recovery replication plan (%s): model is nil. ", id.String())
	}

	return utils.Bool(model.Id != nil), nil
}
