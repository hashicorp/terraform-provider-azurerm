package managementgroup

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccManagementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccManagementGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testManagementGroup_requiresImport(),
				ExpectError: acceptance.RequiresImportError("azurerm_management_group"),
			},
		},
	})
}

func TestAccManagementGroup_nested(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_nested(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists("azurerm_management_group.parent"),
					testCheckManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				ResourceName:      "azurerm_management_group.child",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccManagementGroup_multiLevel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_multiLevel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists("azurerm_management_group.grandparent"),
					testCheckManagementGroupExists("azurerm_management_group.parent"),
					testCheckManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				ResourceName:      "azurerm_management_group.child",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccManagementGroup_multiLevelUpdated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_nested(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists("azurerm_management_group.parent"),
					testCheckManagementGroupExists("azurerm_management_group.child"),
				),
			},
			{
				Config: testManagementGroup_multiLevel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists("azurerm_management_group.grandparent"),
					testCheckManagementGroupExists("azurerm_management_group.parent"),
					testCheckManagementGroupExists("azurerm_management_group.child"),
				),
			},
		},
	})
}

func TestAccManagementGroup_withName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_withName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccManagementGroup_updateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
				),
			},
			{
				Config: testManagementGroup_withName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestmg-%d", data.RandomInteger)),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccManagementGroup_withSubscriptions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group", "test")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_ids.#", "0"),
				),
			},
			{
				Config: testManagementGroup_withSubscriptions(subscriptionID),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_ids.#", "1"),
				),
			},
			{
				Config: testManagementGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagementGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_ids.#", "0"),
				),
			},
		},
	})
}

func testCheckManagementGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ManagementGroups.GroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		groupName := rs.Primary.Attributes["group_id"]

		recurse := false
		resp, err := client.Get(ctx, groupName, "", &recurse, "", "no-cache")
		if err != nil {
			return fmt.Errorf("Bad: Get on managementGroupsClient: %s", err)
		}

		if resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("Management Group does not exist or you do not have proper permissions: %s", groupName)
		}

		return nil
	}
}

func testCheckManagementGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ManagementGroups.GroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_management_group" {
			continue
		}

		name := rs.Primary.Attributes["group_id"]
		recurse := false
		resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")
		if err != nil {
			return nil
		}

		if resp.StatusCode == http.StatusAccepted {
			return fmt.Errorf("Management Group still exists: %s", *resp.Name)
		}
	}

	return nil
}

func testManagementGroup_basic() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
}
`
}

func testManagementGroup_requiresImport() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
}

resource "azurerm_management_group" "import" {
  name = azurerm_management_group.test.name
}
`
}

func testManagementGroup_nested() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "parent" {
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = azurerm_management_group.parent.id
}
`
}

func testManagementGroup_multiLevel() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "grandparent" {
}

resource "azurerm_management_group" "parent" {
  parent_management_group_id = azurerm_management_group.grandparent.id
}

resource "azurerm_management_group" "child" {
  parent_management_group_id = azurerm_management_group.parent.id
}
`
}

func testManagementGroup_withName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name         = "acctestmg-%d"
  display_name = "accTestMG-%d"
}
`, data.RandomInteger, data.RandomInteger)
}

// TODO: switch this out for dynamically creating a subscription once that's supported in the future
func testManagementGroup_withSubscriptions(subscriptionID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  subscription_ids = [
    "%s",
  ]
}
`, subscriptionID)
}
