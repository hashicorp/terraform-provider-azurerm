package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateArmResourceGroupName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "hello",
			ErrCount: 0,
		},
		{
			Value:    "Hello",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "Hello_World",
			ErrCount: 0,
		},
		{
			Value:    "HelloWithNumbers12345",
			ErrCount: 0,
		},
		{
			Value:    "(Did)You(Know)That(Brackets)Are(Allowed)In(Resource)Group(Names)",
			ErrCount: 0,
		},
		{
			Value:    "EndingWithAPeriod.",
			ErrCount: 1,
		},
		{
			Value:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandString(80),
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(81),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateArmResourceGroupName(tc.Value, "azurerm_resource_group")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected validateArmResourceGroupName to trigger '%d' errors for '%s' - got '%d'", tc.ErrCount, tc.Value, len(errors))
		}
	}
}

func init() {
	resource.AddTestSweepers("azurerm_resource_group", &resource.Sweeper{
		Name: "azurerm_resource_group",
		F:    testSweepResourceGroups,
	})
}

func testSweepResourceGroups(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).resourceGroupClient

	log.Printf("Retrieving the Resource Groups..")
	results, err := client.List("", nil)
	if err != nil {
		return fmt.Errorf("Error Listing on Resource Groups: %+v", err)
	}

	for _, profile := range *results.Value {
		if !shouldSweepAcceptanceTestResource(*profile.Name, *profile.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*profile.ID)
		if err != nil {
			return err
		}

		name := resourceId.ResourceGroup

		log.Printf("Deleting Resource Group %q", name)
		deleteResponse, error := client.Delete(name, make(chan struct{}))
		err = <-error
		resp := <-deleteResponse
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return nil
			}

			return err
		}
	}

	return nil
}

func TestAccAzureRMResourceGroup_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMResourceGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists("azurerm_resource_group.test"),
				),
			},
		},
	})
}

func TestAccAzureRMResourceGroup_disappears(t *testing.T) {
	resourceName := "azurerm_resource_group.test"
	ri := acctest.RandInt()
	config := testAccAzureRMResourceGroup_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					testCheckAzureRMResourceGroupDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMResourceGroup_withTags(t *testing.T) {
	resourceName := "azurerm_resource_group.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMResourceGroup_withTags(ri, location)
	postConfig := testAccAzureRMResourceGroup_withTagsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMResourceGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func testCheckAzureRMResourceGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).resourceGroupClient

		resp, err := conn.Get(resourceGroup)
		if err != nil {
			return fmt.Errorf("Bad: Get on resourceGroupClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceGroup := rs.Primary.Attributes["name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).resourceGroupClient

		_, error := conn.Delete(resourceGroup, make(chan struct{}))
		err := <-error
		if err != nil {
			return fmt.Errorf("Bad: Delete on resourceGroupClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMResourceGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).resourceGroupClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_resource_group" {
			continue
		}

		resourceGroup := rs.Primary.ID

		resp, err := conn.Get(resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Resource Group still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMResourceGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}
`, rInt, location)
}

func testAccAzureRMResourceGroup_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`, rInt, location)
}

func testAccAzureRMResourceGroup_withTagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"

    tags {
	environment = "staging"
    }
}
`, rInt, location)
}
