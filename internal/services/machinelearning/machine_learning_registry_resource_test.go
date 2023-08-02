package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningRegistryResource struct{}

func (a MachineLearningRegistryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := registrymanagement.ParseRegistryID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.MachineLearning.RegistryManagementClient.RegistriesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving MachineLearningRegistry %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func TestAccMachineLearningRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, machinelearning.MachineLearningRegistryResource{}.ResourceType(), "test")
	r := MachineLearningRegistryResource{}

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

func TestAccMachineLearningRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, machinelearning.MachineLearningRegistryResource{}.ResourceType(), "test")
	r := MachineLearningRegistryResource{}

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

// update tags still read old value: https://github.com/Azure/azure-rest-api-specs/issues/25200
func TestAccMachineLearningRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, machinelearning.MachineLearningRegistryResource{}.ResourceType(), "test")
	r := MachineLearningRegistryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// data.ImportStep(),
		// {
		// 	Config: r.basicWithNoReplication(data),
		// 	Check: acceptance.ComposeTestCheckFunc(
		// 		check.That(data.ResourceName).ExistsInAzure(r),
		// 	),
		// },
		// data.ImportStep(),
	})
}

func (a MachineLearningRegistryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_machine_learning_registry" "test" {
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlc-%[2]d"

  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = true
  storage_account_type          = "Standard_LRS"

  replication_regions {
    location = "%[4]s"

    storage_account_type = "Standard_LRS"
  }

  tags = {
    key = "example"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (a MachineLearningRegistryResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`




%s

resource "azurerm_machine_learning_registry" "test" {
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlc-%[2]d"

  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = true
  storage_account_type          = "Standard_LRS"

  replication_regions {
    location = "%[4]s"

    storage_account_type = "Standard_LRS"
  }

  tags = {
    key = "example"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Ternary)
}

func (a MachineLearningRegistryResource) basicWithNoReplication(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_machine_learning_registry" "test" {
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlc-%[2]d"

  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = true
  storage_account_type          = "Standard_LRS"

  tags = {
    key = "example"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a MachineLearningRegistryResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_registry" "test" {
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  name                = "accmlc-%[2]d"

  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = true
  storage_account_type          = "Standard_LRS"

  replication_regions {
    location = "%[4]s"

    storage_account_type = "Standard_LRS"
  }

  replication_regions {
    location = "%[5]s"

    storage_account_type        = "Standard_LRS"
    storage_account_hns_enabled = true
  }

  tags = {
    key = "example2"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.Locations.Ternary)
}

func (a MachineLearningRegistryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
