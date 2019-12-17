package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMediaServicesAccount_basic(t *testing.T) {
	resourceName := "azurerm_media_services_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaServicesAccount_basic(ri, rs, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "storage_account.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMediaServicesAccount_multipleAccounts(t *testing.T) {
	resourceName := "azurerm_media_services_account.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaServicesAccount_multipleAccounts(ri, rs, location),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMMediaServicesAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_account.#", "2"),
				),
			},
			{
				Config:   testAccAzureRMMediaServicesAccount_multipleAccountsUpdated(ri, rs, location),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMediaServicesAccount_multiplePrimaries(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaServicesAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMMediaServicesAccount_multiplePrimaries(ri, rs, acceptance.Location()),
				ExpectError: regexp.MustCompile("Only one Storage Account can be set as Primary"),
			},
		},
	})
}

func testCheckAzureRMMediaServicesAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Media service not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Media Services Account: '%s'", name)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on mediaServicesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Media Services Account %q (Resource Group %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMediaServicesAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_media_services_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Media Services Account still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMediaServicesAccount_basic(rInt int, rString, location string) string {
	template := testAccAzureRMMediaServicesAccount_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  storage_account {
    id         = "${azurerm_storage_account.first.id}"
    is_primary = true
  }
}
`, template, rString)
}

func testAccAzureRMMediaServicesAccount_multipleAccounts(rInt int, rString, location string) string {
	template := testAccAzureRMMediaServicesAccount_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  storage_account {
    id         = "${azurerm_storage_account.first.id}"
    is_primary = true
  }

  storage_account {
    id         = "${azurerm_storage_account.second.id}"
    is_primary = false
  }
}
`, template, rString, rString)
}

func testAccAzureRMMediaServicesAccount_multipleAccountsUpdated(rInt int, rString, location string) string {
	template := testAccAzureRMMediaServicesAccount_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  storage_account {
    id         = "${azurerm_storage_account.second.id}"
    is_primary = false
  }

  storage_account {
    id         = "${azurerm_storage_account.first.id}"
    is_primary = true
  }
}
`, template, rString, rString)
}

func testAccAzureRMMediaServicesAccount_multiplePrimaries(rInt int, rString, location string) string {
	template := testAccAzureRMMediaServicesAccount_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "second" {
  name                     = "acctestsa2%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  storage_account {
    id         = "${azurerm_storage_account.first.id}"
    is_primary = true
  }

  storage_account {
    id         = "${azurerm_storage_account.second.id}"
    is_primary = true
  }
}
`, template, rString, rString)
}

func testAccAzureRMMediaServicesAccount_template(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "first" {
  name                     = "acctestsa1%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
`, rInt, location, rString)
}
