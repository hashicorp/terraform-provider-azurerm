package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryHyperVReplicatedVMResource struct{}

func TestAccSiteRecoveryHyperVReplicatedVM_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_replicated_vm", "test")
	r := SiteRecoveryHyperVReplicatedVMResource{}
	hostResource := HyperVHostTestResource{}
	adminPwd := GenerateRandomPassword(10)

	data.ResourceTest(t, r, append(hostResource.PrepareHostTestSteps(data, adminPwd), []acceptance.TestStep{
		{
			PreConfig: func() { time.Sleep(5 * time.Minute) }, // sleep 5 minutes to wait for the host registration fully finished.
			Config:    r.basic(data, adminPwd),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	}...))
}

func TestAccSiteRecoveryHyperVReplicatedVM_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_hyperv_replicated_vm", "test")
	r := SiteRecoveryHyperVReplicatedVMResource{}
	hostResource := HyperVHostTestResource{}
	adminPwd := GenerateRandomPassword(10)

	data.ResourceTest(t, r, append(hostResource.PrepareHostTestSteps(data, adminPwd), []acceptance.TestStep{
		{
			Config: r.complete(data, adminPwd),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	}...))
}

func (SiteRecoveryHyperVReplicatedVMResource) basic(data acceptance.TestData, adminPwd string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "target" {
  name     = "acctest-hyperv-target-%[2]d"
  location = "%[3]s"
}

resource "azurerm_storage_account" "target" {
  name                     = "accttarget%[4]s"
  resource_group_name      = azurerm_resource_group.target.name
  location                 = azurerm_resource_group.target.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "target" {
  name                = "net-%[2]d"
  resource_group_name = azurerm_resource_group.target.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.target.location
}

resource "azurerm_subnet" "target" {
  name                 = "snet-%[2]d"
  resource_group_name  = azurerm_resource_group.target.name
  virtual_network_name = azurerm_virtual_network.target.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_site_recovery_hyperv_replicated_vm" "test" {
  name                      = "acctest-vm-%[2]d"
  hyperv_site_id            = azurerm_site_recovery_services_vault_hyperv_site.test.id
  source_vm_name            = "VM1"
  target_resource_group_id  = azurerm_resource_group.target.id
  target_vm_name            = "target-vm"
  target_storage_account_id = azurerm_storage_account.target.id
  replication_policy_id     = azurerm_site_recovery_hyperv_replication_policy.test.id
  os_type                   = "Windows"
  os_disk_name              = "VM1"
  target_network_id         = azurerm_virtual_network.target.id
  disks_to_include          = ["VM1"]

  network_interface {
    network_name       = "HyperV-NAT"
    target_subnet_name = azurerm_subnet.target.name
  }
}
`, SiteRecoverHyperVReplicationPolicyAssociationResource{}.basic(data, adminPwd), data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (SiteRecoveryHyperVReplicatedVMResource) complete(data acceptance.TestData, adminPwd string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "target" {
  name     = "acctest-hyperv-target-%[2]d"
  location = "%[3]s"
}

resource "azurerm_storage_account" "target" {
  name                     = "accttarget%[4]s"
  resource_group_name      = azurerm_resource_group.target.name
  location                 = azurerm_resource_group.target.location
  account_tier             = "Standard"
  account_kind             = "StorageV2"
  account_replication_type = "LRS"
}

resource "azurerm_virtual_network" "target" {
  name                = "net-%[2]d"
  resource_group_name = azurerm_resource_group.target.name
  address_space       = ["192.168.2.0/24"]
  location            = azurerm_resource_group.target.location
}

resource "azurerm_subnet" "target" {
  name                 = "snet-%[2]d"
  resource_group_name  = azurerm_resource_group.target.name
  virtual_network_name = azurerm_virtual_network.target.name
  address_prefixes     = ["192.168.2.0/24"]
}

resource "azurerm_key_vault" "target" {
  name                        = "kv%[2]d"
  location                    = azurerm_resource_group.target.location
  resource_group_name         = azurerm_resource_group.target.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.target.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
  ]
}

resource "azurerm_key_vault_key" "target" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.target.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_disk_encryption_set" "target" {
  name                = "acctestdes-%[2]d2"
  resource_group_name = azurerm_resource_group.target.name
  location            = azurerm_resource_group.target.location
  key_vault_key_id    = azurerm_key_vault_key.target.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.target.id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
  ]

  tenant_id = azurerm_disk_encryption_set.target.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.target.identity.0.principal_id
}

resource "azurerm_proximity_placement_group" "target" {
  name                = "acctest-replication-%[2]d"
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name
}

resource "azurerm_site_recovery_hyperv_replicated_vm" "test" {
  name                      = "acctest-vm-%[2]d"
  hyperv_site_id            = azurerm_site_recovery_services_vault_hyperv_site.test.id
  source_vm_name            = "VM1"
  target_resource_group_id  = azurerm_resource_group.target.id
  target_vm_name            = "target-vm"
  target_storage_account_id = azurerm_storage_account.target.id
  replication_policy_id     = azurerm_site_recovery_hyperv_replication_policy.test.id
  os_type                   = "Windows"
  os_disk_name              = "VM1"
  target_network_id         = azurerm_virtual_network.target.id
  use_managed_disk_enabled  = true
  managed_disk {
    disk_name                     = "VM1"
    target_disk_encryption_set_id = azurerm_disk_encryption_set.target.id
    target_disk_type              = "Standard_LRS"
  }
  log_storage_account_id             = azurerm_storage_account.target.id
  enable_rdp_or_ssh_on_target_option = "Always"
  network_interface {
    network_name       = "HyperV-NAT"
    target_static_ip   = "10.0.10.5"
    target_subnet_name = azurerm_subnet.target.name
    is_primary         = true
    failover_enabled   = true
  }
  license_type                        = "WindowsServer"
  sql_server_license_type             = "PAYG"
  target_proximity_placement_group_id = azurerm_proximity_placement_group.target.id
  target_vm_tags = {
    tag = "foo"
  }
  target_disk_tags = {
    tag = "foo"
  }
  target_network_interface_tags = {
    tag = "foo"
  }
}
`, SiteRecoverHyperVReplicationPolicyAssociationResource{}.basic(data, adminPwd), data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (t SiteRecoveryHyperVReplicatedVMResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ReplicationProtectedItemsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading site recovery replicated vm (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}
