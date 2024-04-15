package containers_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterTrustedAccessRoleBindingTestResource struct{}

func TestAccKubernetesClusterTrustedAccessRoleBinding_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_trusted_access_role_binding", "test")
	r := KubernetesClusterTrustedAccessRoleBindingTestResource{}

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

func TestAccKubernetesClusterTrustedAccessRoleBinding_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_trusted_access_role_binding", "test")
	r := KubernetesClusterTrustedAccessRoleBindingTestResource{}

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
func (r KubernetesClusterTrustedAccessRoleBindingTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := trustedaccess.ParseTrustedAccessRoleBindingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ContainerService.V20230302Preview.TrustedAccess.RoleBindingsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
func (r KubernetesClusterTrustedAccessRoleBindingTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "test" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  name                  = "acctestkctarb-${var.random_string}"
  roles                 = ["Microsoft.MachineLearningServices/workspaces/mlworkload"]
  source_resource_id    = azurerm_machine_learning_workspace.test.id
}
`, r.template(data))
}

func (r KubernetesClusterTrustedAccessRoleBindingTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "import" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster_trusted_access_role_binding.test.kubernetes_cluster_id
  name                  = azurerm_kubernetes_cluster_trusted_access_role_binding.test.name
  roles                 = azurerm_kubernetes_cluster_trusted_access_role_binding.test.roles
  source_resource_id    = azurerm_kubernetes_cluster_trusted_access_role_binding.test.source_resource_id
}
`, r.basic(data))
}

func (r KubernetesClusterTrustedAccessRoleBindingTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-${var.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}


data "azurerm_client_config" "test" {}


resource "azurerm_key_vault" "test" {
  name                       = "acctest-${var.random_string}"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.test.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}


resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.test.tenant_id
  object_id    = data.azurerm_client_config.test.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}


resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks${var.random_string}"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctestmlw-${var.random_integer}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id
  application_insights_id = azurerm_application_insights.test.id

  identity {
    type = "SystemAssigned"
  }
}


resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}


resource "azurerm_storage_account" "test" {
  name                     = "acctestsa${var.random_string}"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
