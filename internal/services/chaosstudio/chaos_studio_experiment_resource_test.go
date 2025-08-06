package chaosstudio_test

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ChaosStudioExperimentTestResource struct{}

func TestAccChaosStudioExperiment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

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

func TestAccChaosStudioExperiment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

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

func TestAccChaosStudioExperiment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

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

func TestAccChaosStudioExperiment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func TestAccChaosStudioExperiment_exampleAKS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.exampleAKS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccChaosStudioExperiment_multipleSelectors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_chaos_studio_experiment", "test")
	r := ChaosStudioExperimentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleSelectors(data),
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

func (r ChaosStudioExperimentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := experiments.ParseExperimentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ChaosStudio.V20231101.Experiments.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ChaosStudioExperimentTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_experiment" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestcse-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.test.id]
  }

  steps {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.test.urn
        selector_name = "Selector1"
        parameters = {
          abruptShutdown = "false"
        }
        action_type = "continuous"
        duration    = "PT10M"
      }
    }
  }
}
`, r.templateVM(data))
}

func (r ChaosStudioExperimentTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_chaos_studio_experiment" "import" {
  location            = azurerm_chaos_studio_experiment.test.location
  name                = azurerm_chaos_studio_experiment.test.name
  resource_group_name = azurerm_chaos_studio_experiment.test.resource_group_name

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.test.id]
  }

  steps {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.test.urn
        selector_name = "Selector1"
        parameters = {
          abruptShutdown = "false"
        }
        action_type = "continuous"
        duration    = "PT10M"
      }
    }
  }
}
`, r.basic(data))
}

func (r ChaosStudioExperimentTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_experiment" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestcse-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.test.id]
  }

  steps {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.test.urn
        selector_name = "Selector1"
        parameters = {
          abruptShutdown = "false"
        }
        action_type = "continuous"
        duration    = "PT10M"
      }
      actions {
        urn           = azurerm_chaos_studio_capability.test2.urn
        selector_name = "Selector1"
        action_type   = "discrete"
      }
    }
  }
}
`, r.templateVM(data))
}

func (r ChaosStudioExperimentTestResource) exampleAKS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_experiment" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestcse-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.aks.id]
  }

  steps {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.network.urn
        selector_name = "Selector1"
        parameters = {
          jsonSpec = "{\"action\":\"delay\",\"mode\":\"one\",\"selector\":{\"namespaces\":[\"default\"]},\"delay\":{\"latency\":\"200ms\",\"correlation\":\"100\",\"jitter\":\"0ms\"}}}"
        }
        action_type = "discrete"
      }
      actions {
        duration    = "PT10M"
        action_type = "delay"
      }
    }
  }
}
`, r.templateBase(data), r.templateAKS())
}

func (r ChaosStudioExperimentTestResource) multipleSelectors(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

provider "azurerm" {
  features {}
}

resource "azurerm_chaos_studio_experiment" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestcse-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  selectors {
    name                    = "Selector1"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.test.id]
  }
  selectors {
    name                    = "Selector2"
    chaos_studio_target_ids = [azurerm_chaos_studio_target.aks.id]
  }

  steps {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.test.urn
        selector_name = "Selector1"
        parameters = {
          abruptShutdown = "false"
        }
        action_type = "continuous"
        duration    = "PT10M"
      }
      actions {
        urn           = azurerm_chaos_studio_capability.test2.urn
        selector_name = "Selector1"
        action_type   = "discrete"
      }
    }
    branch {
      name = "acctestcse-aks${var.random_string}"
      actions {
        urn           = azurerm_chaos_studio_capability.network.urn
        selector_name = "Selector2"
        parameters = {
          jsonSpec = "{\"action\":\"delay\",\"mode\":\"one\",\"selector\":{\"namespaces\":[\"default\"]},\"delay\":{\"latency\":\"200ms\",\"correlation\":\"100\",\"jitter\":\"0ms\"}}}"
        }
        action_type = "discrete"
      }
    }
  }
}
`, r.templateVM(data), r.templateAKS())
}

func (r ChaosStudioExperimentTestResource) templateVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-${var.random_integer}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-${var.random_integer}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                = "acctestVM-${var.random_integer}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"

  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_chaos_studio_target" "test" {
  location           = azurerm_resource_group.test.location
  target_resource_id = azurerm_linux_virtual_machine.test.id
  target_type        = "Microsoft-VirtualMachine"
}

resource "azurerm_chaos_studio_capability" "test" {
  chaos_studio_target_id = azurerm_chaos_studio_target.test.id
  capability_type        = "Shutdown-1.0"
}

resource "azurerm_chaos_studio_capability" "test2" {
  chaos_studio_target_id = azurerm_chaos_studio_target.test.id
  capability_type        = "Redeploy-1.0"
}
`, r.templateBase(data))
}

func (r ChaosStudioExperimentTestResource) templateAKS() string {
	return `
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

resource "azurerm_chaos_studio_target" "aks" {
  location           = azurerm_resource_group.test.location
  target_resource_id = azurerm_kubernetes_cluster.test.id
  target_type        = "Microsoft-AzureKubernetesServiceChaosMesh"
}

resource "azurerm_chaos_studio_capability" "network" {
  chaos_studio_target_id = azurerm_chaos_studio_target.aks.id
  capability_type        = "NetworkChaos-2.0"
}

resource "azurerm_chaos_studio_capability" "pod" {
  chaos_studio_target_id = azurerm_chaos_studio_target.aks.id
  capability_type        = "PodChaos-2.1"
}
`
}

func (r ChaosStudioExperimentTestResource) templateBase(data acceptance.TestData) string {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctests${var.random_string}"
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
