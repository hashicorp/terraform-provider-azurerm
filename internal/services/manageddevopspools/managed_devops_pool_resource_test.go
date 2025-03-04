package manageddevopspools_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedDevOpsPoolsTestResource struct{}

func TestAccResourceGroupExample_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
    r := ManagedDevOpsPoolsTestResource{}
    
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

func TestAccResourceGroupExample_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
    r := ManagedDevOpsPoolsTestResource{}
    
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

func TestAccResourceGroupExample_complete(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
    r := ManagedDevOpsPoolsTestResource{}
    
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

func TestAccResourceGroupExample_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
    r := ManagedDevOpsPoolsTestResource{}
    
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

func (ManagedDevOpsPoolsTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (ManagedDevOpsPoolsTestResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_managed_devops_pool" "test" {
  name     = "acctest-pool-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagedDevOpsPoolsTestResource) requiresImport(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_managed_devops_pool" "import" {
  name     = azurerm_managed_devops_pool.test.name
  location = azurerm_managed_devops_pool.test.location
}
`, r.basic(data))
}

func (ManagedDevOpsPoolsTestResource) complete(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_managed_devops_pool" "test" {
  name     = "acctest-%d"
  location = "%s"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}