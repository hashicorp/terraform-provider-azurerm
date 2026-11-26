// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/nodepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedhatHcpOpenshiftNodePoolResource struct{}

func TestAccRedhatHcpOpenshiftNodePool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("replicas").HasValue("1"),
				check.That(data.ResourceName).Key("platform.0.vm_size").HasValue("Standard_D4s_v3"),
				check.That(data.ResourceName).Key("platform.0.subnet_id").MatchesOtherKey(check.That("azurerm_subnet.worker_subnet").Key("id")),
				check.That(data.ResourceName).Key("hcp_openshift_cluster_id").MatchesOtherKey(check.That("azurerm_redhat_hcp_openshift_cluster.test").Key("id")),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftNodePool_platform(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.platform(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("platform.0.vm_size").HasValue("Standard_D4s_v5"),
				check.That(data.ResourceName).Key("platform.0.subnet_id").MatchesOtherKey(check.That("azurerm_subnet.worker_subnet").Key("id")),
				check.That(data.ResourceName).Key("platform.0.availability_zone").HasValue("1"),
				check.That(data.ResourceName).Key("platform.0.enable_encryption_at_host").HasValue("true"),
				check.That(data.ResourceName).Key("platform.0.os_disk.0.disk_storage_account_type").HasValue("Premium_LRS"),
				check.That(data.ResourceName).Key("platform.0.os_disk.0.disk_type").HasValue("Managed"),
				check.That(data.ResourceName).Key("platform.0.os_disk.0.encryption_set_id").MatchesOtherKey(check.That("azurerm_disk_encryption_set.test").Key("id")),
				check.That(data.ResourceName).Key("platform.0.os_disk.0.size_gb").HasValue("128"),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeFromState(r.platformDependencies(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftNodePool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
		{
			Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftNodePool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("replicas").HasValue("2"),
				check.That(data.ResourceName).Key("node_drain_timeout_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("labels.workload").HasValue("acceptance"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acceptance"),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftNodePool_autoscaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
		{
			Config: r.autoscaling(data, 2, 4),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_scaling.0.min").HasValue("2"),
				check.That(data.ResourceName).Key("auto_scaling.0.max").HasValue("4"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoscaling(data, 2, 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_scaling.0.max").HasValue("5"),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	})
}

func TestAccRedhatHcpOpenshiftNodePool_invalidAutoscaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.autoscaling(data, 5, 2),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("auto_scaling.min.*must not exceed.*auto_scaling.max"),
		},
	})
}

func (r RedhatHcpOpenshiftNodePoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nodepools.ParseNodePoolID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.RedHatOpenShift.NodePoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r RedhatHcpOpenshiftNodePoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_node_pool" "test" {
  name                     = "acctestnp%[2]s"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id
  replicas                 = 1

  platform {
    vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }

  version {
    id = "4.22.20"
  }
}
`, RedhatHcpOpenshiftClusterResource{}.basic(data), data.RandomString)
}

func (r RedhatHcpOpenshiftNodePoolResource) platformDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "node_pool" {
  name                       = "acctestnpkv%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  purge_protection_enabled   = true
  rbac_authorization_enabled = true
}

resource "azurerm_role_assignment" "node_pool_key_vault_crypto_officer" {
  scope                = azurerm_key_vault.node_pool.id
  role_definition_name = "Key Vault Crypto Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "node_pool" {
  name         = "node-pool-disk-encryption-key"
  key_vault_id = azurerm_key_vault.node_pool.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "unwrapKey", "wrapKey"]

  depends_on = [azurerm_role_assignment.node_pool_key_vault_crypto_officer]
}

resource "azurerm_user_assigned_identity" "disk_encryption_set" {
  name                = "acctestnpdes-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "disk_encryption_set_key_vault" {
  scope                            = azurerm_key_vault.node_pool.id
  role_definition_name             = "Key Vault Crypto Service Encryption User"
  principal_id                     = azurerm_user_assigned_identity.disk_encryption_set.principal_id
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "disk_encryption_set_capi_reader" {
  scope                            = azurerm_disk_encryption_set.test.id
  principal_id                     = local.control_plane_operators["cluster-api-azure"].principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_role_assignment" "disk_encryption_set_service_reader" {
  scope                            = azurerm_disk_encryption_set.test.id
  principal_id                     = azurerm_user_assigned_identity.service.principal_id
  role_definition_id               = local.role_definition_ids.reader
  principal_type                   = "ServicePrincipal"
  skip_service_principal_aad_check = true
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestnpdes-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  key_vault_key_id    = azurerm_key_vault_key.node_pool.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.disk_encryption_set.id]
  }

  depends_on = [azurerm_role_assignment.disk_encryption_set_key_vault]
}
`, RedhatHcpOpenshiftClusterResource{}.basicWithDependencies(data, `
		# The node pool is removed from state without destruction, leaving its disks
		# until cluster deletion. Delete the cluster before the disk encryption set
		# so those disks no longer reference it.
		azurerm_disk_encryption_set.test,
`), data.RandomString)
}

func (r RedhatHcpOpenshiftNodePoolResource) platform(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_redhat_hcp_openshift_node_pool" "test" {
  name                     = "platform"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id
  replicas                 = 1

  platform {
    vm_size                   = "Standard_D4s_v5"
    subnet_id                 = azurerm_subnet.worker_subnet.id
    availability_zone         = "1"
    enable_encryption_at_host = true

    os_disk {
      disk_storage_account_type = "Premium_LRS"
      disk_type                 = "Managed"
      encryption_set_id         = azurerm_disk_encryption_set.test.id
      size_gb                   = 128
    }
  }

  version {
    id = "4.22.20"
  }

  depends_on = [
    azurerm_role_assignment.disk_encryption_set_capi_reader,
    azurerm_role_assignment.disk_encryption_set_service_reader,
  ]
}
`, r.platformDependencies(data))
}

// A bug in the HCP service currently prevents proper destruction of node pool resources
// because it won't allow deletion of a cluster's last node pool once at least one has been created.
// Once the bug is fixed, we can remove the `removed` block workaround.
func (r RedhatHcpOpenshiftNodePoolResource) removeFromState(dependencies string) string {
	return fmt.Sprintf(`
%s

terraform {
  required_version = ">= 1.7.0"
}

removed {
  from = azurerm_redhat_hcp_openshift_node_pool.test

  lifecycle {
    destroy = false
  }
}
`, dependencies)
}

func (r RedhatHcpOpenshiftNodePoolResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_node_pool" "test" {
  name                       = "acctestnp%[2]s"
  hcp_openshift_cluster_id   = azurerm_redhat_hcp_openshift_cluster.test.id
  node_drain_timeout_minutes = 30
  replicas                   = 2

  platform {
    vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }

  version {
    id = "4.22.20"
  }

  labels = {
    workload = "acceptance"
  }

  tags = {
    environment = "acceptance"
  }
}
`, RedhatHcpOpenshiftClusterResource{}.basic(data), data.RandomString)
}

func (r RedhatHcpOpenshiftNodePoolResource) autoscaling(data acceptance.TestData, min, max int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_node_pool" "test" {
  name                     = "acctestnp%[2]s"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id

  platform {
    vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }

  auto_scaling {
    min = %[3]d
    max = %[4]d
  }

  version {
    id = "4.22.20"
  }
}
`, RedhatHcpOpenshiftClusterResource{}.basic(data), data.RandomString, min, max)
}

func (r RedhatHcpOpenshiftNodePoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_redhat_hcp_openshift_node_pool" "test" {
  name                     = "acctestnp%[2]s"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.test.id
  replicas                 = 2

  platform {
    vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }

  version {
    id = "4.22.20"
  }
}

resource "azurerm_redhat_hcp_openshift_node_pool" "import" {
  name                     = azurerm_redhat_hcp_openshift_node_pool.test.name
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_node_pool.test.hcp_openshift_cluster_id
  replicas                 = 1

  platform {
    vm_size = azurerm_redhat_hcp_openshift_node_pool.test.platform[0].vm_size
  }

  version {
    id = azurerm_redhat_hcp_openshift_node_pool.test.version[0].id
  }
}
`, RedhatHcpOpenshiftClusterResource{}.basic(data), data.RandomString)
}
