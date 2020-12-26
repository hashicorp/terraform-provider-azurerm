package resource_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMManagementLock_resourceGroupReadOnlyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMManagementLock_requiresImport),
		},
	})
}

func TestAccAzureRMManagementLock_resourceGroupReadOnlyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupReadOnlyComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_publicIPReadOnlyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_publicIPReadOnlyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_publicIPCanNotDeleteBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_publicIPCanNotDeleteBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_subscriptionReadOnlyBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_subscriptionReadOnlyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMManagementLockExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.LocksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		resp, err := client.GetByScope(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Management Lock %q (Scope %q) does not exist", name, scope)
			}

			return fmt.Errorf("Bad: Get on managementLocksClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMManagementLockDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.LocksClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_management_lock" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		resp, err := client.GetByScope(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagementLock_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "import" {
  name       = azurerm_management_lock.test.name
  scope      = azurerm_management_lock.test.scope
  lock_level = azurerm_management_lock.test.lock_level
}
`, template)
}

func testAccAzureRMManagementLock_resourceGroupReadOnlyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "ReadOnly"
  notes      = "Hello, World!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "CanNotDelete"
  notes      = "Hello, World!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMManagementLock_publicIPReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_public_ip.test.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagementLock_publicIPCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_public_ip.test.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMManagementLock_subscriptionReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = data.azurerm_subscription.current.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger)
}

func testAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = data.azurerm_subscription.current.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger)
}
