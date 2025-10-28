// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcKubernetesProvisionedClusterResource struct{}

func TestAccArcKubernetesProvisionedCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster", "test")
	r := ArcKubernetesProvisionedClusterResource{}
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

func TestAccArcKubernetesProvisionedCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster", "test")
	r := ArcKubernetesProvisionedClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccArcKubernetesProvisionedCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster", "test")
	r := ArcKubernetesProvisionedClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccArcKubernetesProvisionedCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_arc_kubernetes_provisioned_cluster", "test")
	r := ArcKubernetesProvisionedClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ArcKubernetesProvisionedClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := arckubernetes.ParseConnectedClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ArcKubernetes.ArcKubernetesClient.ConnectedClusterGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ArcKubernetesProvisionedClusterResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_arc_kubernetes_provisioned_cluster" "test" {
  name                = "acctest-akpc-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesProvisionedClusterResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_arc_kubernetes_provisioned_cluster" "import" {
  name                = azurerm_arc_kubernetes_provisioned_cluster.test.name
  resource_group_name = azurerm_arc_kubernetes_provisioned_cluster.test.resource_group_name
  location            = azurerm_arc_kubernetes_provisioned_cluster.test.location
  identity {
    type = "SystemAssigned"
  }
}
`, config)
}

func (r ArcKubernetesProvisionedClusterResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_arc_kubernetes_provisioned_cluster" "test" {
  name                           = "acctest-akpc-%[2]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  arc_agent_auto_upgrade_enabled = false
  identity {
    type = "SystemAssigned"
  }

  azure_active_directory {
    azure_rbac_enabled = false
  }

  tags = {
    ENV = "TestUpdate"
  }
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesProvisionedClusterResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctestADG-arck8s-%[2]d"
  owners           = [data.azurerm_client_config.current.object_id]
  security_enabled = true
}

resource "azurerm_arc_kubernetes_provisioned_cluster" "test" {
  name                           = "acctest-akpc-%[2]d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  arc_agent_auto_upgrade_enabled = true
  arc_agent_desired_version      = "1.18.3"
  identity {
    type = "SystemAssigned"
  }

  azure_active_directory {
    azure_rbac_enabled     = true
    admin_group_object_ids = [azuread_group.test.object_id]
    tenant_id              = data.azurerm_client_config.current.tenant_id
  }

  tags = {
    ENV = "Test"
    FOO = "BAR"
  }
}
`, template, data.RandomInteger)
}

func (r ArcKubernetesProvisionedClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-akpc-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
