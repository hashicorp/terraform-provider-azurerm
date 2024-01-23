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
func (r ChaosStudioExperimentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := experiments.ParseExperimentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ChaosStudio.V20231101.Experiments.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
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
  selector {
    id   = "val-${var.random_string}"
    type = "List"
  }
  step {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      action {
        name = "acctestcse-${var.random_string}"
        type = "val-${var.random_string}"
      }
    }
  }
}
`, r.template(data))
}

func (r ChaosStudioExperimentTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_chaos_studio_experiment" "import" {
  location            = azurerm_chaos_studio_experiment.test.location
  name                = azurerm_chaos_studio_experiment.test.name
  resource_group_name = azurerm_chaos_studio_experiment.test.resource_group_name
  selector            = azurerm_chaos_studio_experiment.test.selector
  step                = azurerm_chaos_studio_experiment.test.step
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
  selector {
    id   = "val-${var.random_string}"
    type = "List"
    filter {
      type = "Simple"
    }
  }
  step {
    name = "acctestcse-${var.random_string}"
    branch {
      name = "acctestcse-${var.random_string}"
      action {
        name = "acctestcse-${var.random_string}"
        type = "val-${var.random_string}"
      }
    }
  }
  tags = {
    environment = "terraform-acctests"
    some_key    = "some-value"
  }
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
}
`, r.template(data))
}

func (r ChaosStudioExperimentTestResource) template(data acceptance.TestData) string {
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

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks${var.random_string}"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_chaos_studio_target" "test" {
  location           = azurerm_resource_group.test.location
  target_resource_id = azurerm_kubernetes_cluster.test.id
  target_type        = "Microsoft-AzureKubernetesServiceChaosMesh"
}

resource "azurerm_chaos_studio_capability" "test" {
  chaos_studio_target_id = azurerm_chaos_studio_target.test.id
  capability_type        = "NetworkChaos-2.0"
}

resource "azurerm_chaos_studio_capability" "test" {
  location           = azurerm_resource_group.test.location
  target_resource_id = azurerm_storage_account.test.id
  target_type        = "PodChaos-2.1"
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
