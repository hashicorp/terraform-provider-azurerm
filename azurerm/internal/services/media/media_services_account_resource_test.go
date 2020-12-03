package media_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
)

func TestAccMediaServicesAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaServicesAccount_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccMediaServicesAccount_multipleAccounts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaServicesAccount_multipleAccounts(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckMediaServicesAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account.#", "2"),
				),
			},
			{
				Config:   testAccMediaServicesAccount_multipleAccountsUpdated(data),
				PlanOnly: true,
			},
			data.ImportStep(),
		},
	})
}

func TestAccMediaServicesAccount_multiplePrimaries(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccMediaServicesAccount_multiplePrimaries(data),
				ExpectError: regexp.MustCompile("Only one Storage Account can be set as Primary"),
			},
		},
	})
}

func TestAccMediaServicesAccount_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_services_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaServicesAccount_identitySystemAssigned(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckMediaServicesAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Media service not found: %s", resourceName)
		}

		id, err := parse.MediaServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on mediaServicesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Media Services Account %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckMediaServicesAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_media_services_account" {
			continue
		}

		id, err := parse.MediaServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Media Services Account still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccMediaServicesAccount_basic(data acceptance.TestData) string {
	template := testAccMediaServicesAccount_template(data)
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

func testAccMediaServicesAccount_multipleAccounts(data acceptance.TestData) string {
	template := testAccMediaServicesAccount_template(data)
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

func testAccMediaServicesAccount_multipleAccountsUpdated(data acceptance.TestData) string {
	template := testAccMediaServicesAccount_template(data)
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

func testAccMediaServicesAccount_multiplePrimaries(data acceptance.TestData) string {
	template := testAccMediaServicesAccount_template(data)
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

func testAccMediaServicesAccount_identitySystemAssigned(data acceptance.TestData) string {
	template := testAccMediaServicesAccount_template(data)
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

func testAccMediaServicesAccount_template(data acceptance.TestData) string {
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
