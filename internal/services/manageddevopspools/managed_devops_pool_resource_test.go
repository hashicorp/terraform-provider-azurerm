package manageddevopspools_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedDevOpsPoolResource struct{}

func TestAccManagedDevOpsPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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
	})
}

func (ManagedDevOpsPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pools.ParsePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ManagedDevOpsPools.PoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ManagedDevOpsPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_devops_pool" "test" {
  name     = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  organization_profile {
    kind = "AzureDevOps"
    organizations {
      parallelism = 1
      url         = "https://dev.azure.com/managed-org-demo"
    }
  }

  agent_profile {
    kind = "Stateless"
  }

  fabric_profile {
    kind = "Vmss"
    images {
      resource_id = "/Subscriptions/e21f7bce-1728-44ad-a62d-344064a0d69a/Providers/Microsoft.Compute/Locations/australiaeast/publishers/canonical/artifacttypes/vmimage/offers/0001-com-ubuntu-server-focal/skus/20_04-lts-gen2/versions/latest"
      buffer = "*"
    }
    sku {
      name = "Standard_D2ads_v5"
    }
  }
}
`, r.template(data))
}

func (r ManagedDevOpsPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_devops_pool" "import" {
  name     			  = azurerm_managed_devops_pool.test.name
  location            = azurerm_managed_devops_pool.test.location
  resource_group_name = azurerm_managed_devops_pool.test.resource_group_name

  maximum_concurrency = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  organization_profile {
    kind = "AzureDevOps"
    organizations {
      parallelism = 1
      url         = "https://dev.azure.com/managed-org-demo"
    }
  }

  agent_profile {
    kind = "Stateless"
  }

  fabric_profile {
    kind = "Vmss"
    images {
      resource_id = "/Subscriptions/e21f7bce-1728-44ad-a62d-344064a0d69a/Providers/Microsoft.Compute/Locations/australiaeast/publishers/canonical/artifacttypes/vmimage/offers/0001-com-ubuntu-server-focal/skus/20_04-lts-gen2/versions/latest"
      buffer = "*"
    }
    sku {
      name = "Standard_D2ads_v5"
    }
  }
}
`, r.basic(data))
}

func (r ManagedDevOpsPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

esource "azurerm_managed_devops_pool" "test" {
  name     = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  organization_profile {
    kind = "AzureDevOps"
    organizations {
      parallelism = 1
      url         = "https://dev.azure.com/managed-org-demo"
    }
  }

  agent_profile {
    kind = "Stateless"
	resource_predictions {
      time_zone = "UTC"
      days_data = "[{},{\"09:00:00\":1,\"17:00:00\": 0},{},{},{},{},{}]"
    }
  }

  fabric_profile {
    kind = "Vmss"
    images {
      resource_id = "/Subscriptions/e21f7bce-1728-44ad-a62d-344064a0d69a/Providers/Microsoft.Compute/Locations/australiaeast/publishers/canonical/artifacttypes/vmimage/offers/0001-com-ubuntu-server-focal/skus/20_04-lts-gen2/versions/latest"
      buffer = "*"
    }
    sku {
      name = "Standard_D2ads_v5"
    }
  }

  tags = {
    Environment = "ppe"
    Project     = "Terraform"
  }
}
`, r.template(data))
}

func (ManagedDevOpsPoolResource) template(data acceptance.TestData) string {
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
  name                = "acctestuami-${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_center" "test" {
  name                = "acctestdc-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_dev_center_project" "test" {
  name                = "acctestproj-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
