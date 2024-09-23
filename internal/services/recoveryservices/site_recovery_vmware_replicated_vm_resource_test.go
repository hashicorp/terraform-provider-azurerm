// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryVMWareReplicatedVmResource struct {
	SubscriptionId   string
	VaultName        string
	VaultRgName      string
	SourceVMName     string
	ApplianceName    string
	Location         string
	Credential       string
	SourceMacAddress string
}

func (r SiteRecoveryVMWareReplicatedVmResource) preCheck(t *testing.T) {
	if r.SubscriptionId == "" {
		t.Skipf("subscription id is empty")
	}
	if r.VaultName == "" {
		t.Skipf("vault name is empty")
	}
	if r.VaultRgName == "" {
		t.Skipf("vault resource group name is empty")
	}
	if r.SourceVMName == "" {
		t.Skipf("`ARM_TEST_VMWARE_SOURCE_VM_NAME` must be set for acceptance tests!")
	}
	if r.ApplianceName == "" {
		t.Skipf("`ARM_TEST_VMWARE_APPLIANCE_NAME` must be set for acceptance tests!")
	}
	if r.Location == "" {
		t.Skipf("`ARM_TEST_VMWARE_VAULT_LOCATION` must be set for acceptance tests!")
	}
	if r.Credential == "" {
		t.Skipf("`ARM_TEST_VMWARE_CREDENTIAL_NAME` must be set for acceptance tests!")
	}
	if r.SourceMacAddress == "" {
		t.Skipf("`ARM_TEST_VMWARE_SOURCE_MAC_ADDRESS` must be set for acceptance tests!")
	}
}

func newSiteRecoveryVMWareReplicatedVMResource(vaultId, sourceVMName, applianceName, location, credential, sourceMacAddress string) (SiteRecoveryVMWareReplicatedVmResource, error) {
	parsedVaultId, err := replicationprotecteditems.ParseVaultID(vaultId)
	if err != nil {
		return SiteRecoveryVMWareReplicatedVmResource{}, fmt.Errorf("parsing %q: %+v", vaultId, err)
	}
	return SiteRecoveryVMWareReplicatedVmResource{
		SubscriptionId:   parsedVaultId.SubscriptionId,
		VaultName:        parsedVaultId.VaultName,
		VaultRgName:      parsedVaultId.ResourceGroupName,
		SourceVMName:     sourceVMName,
		ApplianceName:    applianceName,
		Location:         location,
		Credential:       credential,
		SourceMacAddress: sourceMacAddress,
	}, nil
}

func (r SiteRecoveryVMWareReplicatedVmResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationProtectedItemsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func TestAccSiteVMWareRecoveryReplicatedVM_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_vmware_replicated_vm", "test")
	r, err := newSiteRecoveryVMWareReplicatedVMResource(
		os.Getenv("ARM_TEST_VMWARE_VAULT_ID"),
		os.Getenv("ARM_TEST_VMWARE_SOURCE_VM_NAME"),
		os.Getenv("ARM_TEST_VMWARE_APPLIANCE_NAME"),
		os.Getenv("ARM_TEST_VMWARE_VAULT_LOCATION"),
		os.Getenv("ARM_TEST_VMWARE_CREDENTIAL_NAME"),
		os.Getenv("ARM_TEST_VMWARE_SOURCE_MAC_ADDRESS"),
	)
	if err != nil {
		t.Skipf("failed to create SiteRecoveryVMWareReplicatedVmResource: %+v", err)
	}

	r.preCheck(t)

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

func (r SiteRecoveryVMWareReplicatedVmResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  subscription_id = "%[1]s"
  features {
  }
}

resource "azurerm_resource_group" "target" {
  name     = "acctestRG-%[4]d"
  location = "%[6]s"
}

resource "azurerm_storage_account" "target" {
  name                     = "acct%[5]s"
  resource_group_name      = azurerm_resource_group.target.name
  location                 = azurerm_resource_group.target.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "target" {
  name                = "acctestvn-%[4]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name
}

resource "azurerm_subnet" "target" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.target.name
  virtual_network_name = azurerm_virtual_network.target.name
  address_prefixes     = ["10.0.2.0/24"]
}

data "azurerm_recovery_services_vault" "vault" {
  name                = "%[2]s"
  resource_group_name = "%[3]s"
}

resource "azurerm_site_recovery_vmware_replication_policy" "test" {
  recovery_vault_id                                    = data.azurerm_recovery_services_vault.vault.id
  name                                                 = "acctest-policy-%[4]d"
  recovery_point_retention_in_minutes                  = 1440
  application_consistent_snapshot_frequency_in_minutes = 240
}

resource "azurerm_site_recovery_vmware_replication_policy_association" "test" {
  name              = "acctest-%[4]d"
  recovery_vault_id = data.azurerm_recovery_services_vault.vault.id
  policy_id         = azurerm_site_recovery_vmware_replication_policy.test.id
}

resource "azurerm_site_recovery_vmware_replicated_vm" "test" {
  name                                       = "acct%[4]d"
  recovery_vault_id                          = data.azurerm_recovery_services_vault.vault.id
  source_vm_name                             = "%[7]s"
  appliance_name                             = "%[8]s"
  recovery_replication_policy_id             = azurerm_site_recovery_vmware_replication_policy_association.test.policy_id
  physical_server_credential_name            = "%[9]s"
  license_type                               = "NotSpecified"
  target_boot_diagnostics_storage_account_id = azurerm_storage_account.target.id
  target_vm_name                             = "%[7]s"
  target_resource_group_id                   = azurerm_resource_group.target.id
  default_log_storage_account_id             = azurerm_storage_account.target.id
  default_recovery_disk_type                 = "Standard_LRS"
  target_network_id                          = azurerm_virtual_network.target.id

  network_interface {
    source_mac_address = "%[10]s"
    target_subnet_name = azurerm_subnet.target.name
    is_primary         = true
  }

  lifecycle {
    ignore_changes = [
      target_vm_size,
      default_log_storage_account_id,
      default_recovery_disk_type,
      managed_disk,
      network_interface
    ]
  }
}
`, r.SubscriptionId, r.VaultName, r.VaultRgName, data.RandomInteger, data.RandomString, r.Location, r.SourceVMName, r.ApplianceName, r.Credential, r.SourceMacAddress)
}
