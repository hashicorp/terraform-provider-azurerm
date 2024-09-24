// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ComputeClusterResource struct{}

func TestAccComputeCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("scale_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("scale_settings.0.max_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.min_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.scale_down_nodes_after_idle_duration").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("scale_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("scale_settings.0.max_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.min_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.scale_down_nodes_after_idle_duration").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeCluster_recreateVmSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.recreateVmSize(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("scale_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("scale_settings.0.max_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.min_node_count").Exists(),
				check.That(data.ResourceName).Key("scale_settings.0.scale_down_nodes_after_idle_duration").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccComputeCluster_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_cluster", "test")
	r := ComputeClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r ComputeClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	computeClusterClient := client.MachineLearning.MachineLearningComputes
	id, err := machinelearningcomputes.ParseComputeID(state.ID)
	if err != nil {
		return nil, err
	}

	computeResource, err := computeClusterClient.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Compute Cluster %q: %+v", state.ID, err)
	}
	return utils.Bool(computeResource.Model.Properties != nil), nil
}

func (r ComputeClusterResource) basic(data acceptance.TestData) string {
	template := r.template_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_DS2_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) complete(data acceptance.TestData) string {
	template := r.template_complete(data)
	return fmt.Sprintf(`
%s
variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_DS2_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  subnet_resource_id            = azurerm_subnet.test.id
  node_public_ip_enabled        = false
  description                   = "Machine Learning"
  tags = {
    environment = "test"
  }
  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned"
  }

  ssh_public_access_enabled = false
  ssh {
    admin_username = "adminuser"
    key_value      = var.ssh_key
  }
  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_private_endpoint.test,
  ]
}
`, template, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) recreateVmSize(data acceptance.TestData) string {
	template := r.template_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_D1_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  local_auth_enabled            = false

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_cluster" "import" {
  name                          = azurerm_machine_learning_compute_cluster.test.name
  location                      = azurerm_machine_learning_compute_cluster.test.location
  vm_priority                   = azurerm_machine_learning_compute_cluster.test.vm_priority
  vm_size                       = azurerm_machine_learning_compute_cluster.test.vm_size
  machine_learning_workspace_id = azurerm_machine_learning_compute_cluster.test.machine_learning_workspace_id

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT2M" # 120 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func (r ComputeClusterResource) identitySystemAssigned(data acceptance.TestData) string {
	template := r.template_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_DS2_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "CC-%d"
  location                      = azurerm_resource_group.test.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_DS2_V2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	template := r.template_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_compute_cluster" "test" {
  name        = "CC-%d"
  location    = azurerm_resource_group.test.location
  vm_priority = "LowPriority"
  vm_size     = "STANDARD_DS2_V2"

  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8))
}

func (r ComputeClusterResource) template_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
  tags = {
    "stage" = "test"
  }
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acckv%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW%[5]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary,
		data.RandomInteger, data.RandomIntOfLength(15), data.RandomIntOfLength(16),
		data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ComputeClusterResource) template_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
  tags = {
    "stage" = "test"
  }
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acckv%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acc-MLW%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.api.azureml.ms"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "test-vlink"
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
}

resource "azurerm_private_endpoint" "test" {
  name                = "test-pe-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "test-mlworkspace-%[1]d"
    private_connection_resource_id = azurerm_machine_learning_workspace.test.id
    subresource_names              = ["amlworkspace"]
    is_manual_connection           = false
  }

  private_dns_zone_group {
    name                 = "test"
    private_dns_zone_ids = [azurerm_private_dns_zone.test.id]
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "accvnet%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "accsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "test-nsg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "29876-29877"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
