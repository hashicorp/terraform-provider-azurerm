package azurerm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMDataFactoryV2_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDataFactoryV2_tags(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.environment", "production"),
				),
			},
		},
	})
}

func TestAccAzureRMDataFactoryV2_importWithTags(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_tags(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "azurerm_data_factory_v2.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMDataFactoryV2_tagsUpdated(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_tags(ri, testLocation())
	updatedConfig := testAccAzureRMDataFactoryV2_tagsUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.environment", "production"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.environment", "production"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "tags.updated", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMDataFactoryV2_identity(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_identity(ri, testLocation())
	match := regexp.MustCompile("^(\\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\\}{0,1})$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "identity.#", "1"),
					resource.TestCheckResourceAttr("azurerm_data_factory_v2.test", "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr("azurerm_data_factory_v2.test", "identity.0.principal_id", match),
					resource.TestMatchResourceAttr("azurerm_data_factory_v2.test", "identity.0.tenant_id", match),
				),
			},
		},
	})
}

func TestAccAzureRMDataFactoryV2_disappears(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDataFactoryV2_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryV2Exists("azurerm_data_factory_v2.test"),
					testCheckAzureRMDataFactoryV2Disappears("azurerm_data_factory_v2.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMDataFactoryV2Exists(name string) resource.TestCheckFunc {
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

func testCheckAzureRMDataFactoryV2Disappears(name string) resource.TestCheckFunc {
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

func testCheckAzureRMDataFactoryV2Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).dataFactoryClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_v2" {
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

func testAccAzureRMDataFactoryV2_basic(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory_v2" "test" {
    name                = "acctestdfv2%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactoryV2_tags(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory_v2" "test" {
    name                = "acctestdfv2%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    tags {
      environment = "production"
    }
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactoryV2_tagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory_v2" "test" {
    name                = "acctestdfv2%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    tags {
      environment = "production"
      updated     = "true"
    }
  }
`, rInt, location, rInt)
}

func testAccAzureRMDataFactoryV2_identity(rInt int, location string) string {
	return fmt.Sprintf(`
  resource "azurerm_resource_group" "test" {
    name     = "acctestrg-%d"
    location = "%s"
  }
  resource "azurerm_data_factory_v2" "test" {
    name                = "acctestdfv2%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    
    identity {
      type = "SystemAssigned"
    }
  }
`, rInt, location, rInt)
}
