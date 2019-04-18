package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMDataFactory_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_basic(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
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

func TestAccAzureRMDataFactory_tags(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_tags(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
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

func TestAccAzureRMDataFactory_tagsUpdated(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_tags(ri, testLocation())
	updatedConfig := testAccAzureRMDataFactory_tagsUpdated(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(resourceName, "tags.updated", "true"),
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

func TestAccAzureRMDataFactory_identity(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_identity(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "identity.#"),
					resource.TestCheckResourceAttrSet(resourceName, "identity.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(resourceName, "identity.0.tenant_id"),
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

func TestAccAzureRMDataFactory_disappears(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_basic(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					testCheckAzureRMDataFactoryDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMDataFactory_github(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactory_github(ri, testLocation())
	config2 := testAccAzureRMDataFactory_githubUpdated(ri, testLocation())
	resourceName := "azurerm_data_factory.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.account_name", fmt.Sprintf("acctestrg-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.git_url", "https://github.com/terraform-providers/"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.repository_name", "terraform-provider-azurerm"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.branch_name", "master"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.root_folder", "/"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.account_name", fmt.Sprintf("acctestrg-%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.git_url", "https://github.com/terraform-providers/"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.repository_name", "terraform-provider-azuread"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.branch_name", "stable-website"),
					resource.TestCheckResourceAttr(resourceName, "github_configuration.0.root_folder", "/azuread"),
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

func testCheckAzureRMDataFactoryExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactoryClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Delete on dataFactoryClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory still exists:\n%#v", resp.FactoryProperties)
		}
	}

	return nil
}

func testAccAzureRMDataFactory_basic(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }

  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactory_tags(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }

  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    tags = {
      environment = "production"
    }
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactory_tagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }

  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    tags = {
      environment = "production"
      updated     = "true"
    }
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactory_identity(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }

  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    
    identity {
      type = "SystemAssigned"
    }
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactory_github(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }

  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
		
    github_configuration {
      git_url         = "https://github.com/terraform-providers/"
      repository_name = "terraform-provider-azurerm"
      branch_name     = "master"
      root_folder     = "/"
      account_name    = "acctestrg-%d"
    }
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDataFactory_githubUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory" "test" {
    name                = "acctestdf%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
		
    github_configuration {
      git_url         = "https://github.com/terraform-providers/"
      repository_name = "terraform-provider-azuread"
      branch_name     = "stable-website"
      root_folder     = "/azuread"
      account_name    = "acctestrg-%d"
    }
  }
`, rInt, location, rInt, rInt)
}
