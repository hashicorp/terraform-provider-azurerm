package media_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MediaServicesAccountResource struct {
}

func TestAccMediaServicesAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account.#").HasValue("1"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMediaServicesAccount_multipleAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleAccounts(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config:   r.multipleAccountsUpdated(data),
			PlanOnly: true,
		},
		data.ImportStep(),
	})
}

func TestAccMediaServicesAccount_multiplePrimaries(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.multiplePrimaries(data),
			ExpectError: regexp.MustCompile("Only one Storage Account can be set as Primary"),
		},
	})
}

func TestAccMediaServicesAccount_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")
	r := MediaServicesAccountResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func (MediaServicesAccountResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.MediaServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.ServicesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Media Services Account %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ServiceProperties != nil), nil
}

func (r MediaServicesAccountResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}
`, template, data.RandomString)
}

func (r MediaServicesAccountResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "import" {
  name                = azurerm_media_services_account.test.name
  location            = azurerm_media_services_account.test.location
  resource_group_name = azurerm_media_services_account.test.resource_group_name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  tags = {
    environment = "staging"
  }
}
`, template)
}

func (MediaServicesAccountResource) multipleAccounts(data acceptance.TestData) string {
	template := MediaServicesAccountResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = false
  }
}
`, template, data.RandomString, data.RandomString)
}

func (MediaServicesAccountResource) multipleAccountsUpdated(data acceptance.TestData) string {
	template := MediaServicesAccountResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = false
  }

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }
}
`, template, data.RandomString, data.RandomString)
}

func (MediaServicesAccountResource) multiplePrimaries(data acceptance.TestData) string {
	template := MediaServicesAccountResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  storage_account {
    id         = azurerm_storage_account.second.id
    is_primary = true
  }
}
`, template, data.RandomString, data.RandomString)
}

func (MediaServicesAccountResource) identitySystemAssigned(data acceptance.TestData) string {
	template := MediaServicesAccountResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.first.id
    is_primary = true
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomString)
}

func (MediaServicesAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "first" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
