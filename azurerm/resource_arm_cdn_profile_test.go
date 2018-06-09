package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("azurerm_cdn_profile", &resource.Sweeper{
		Name: "azurerm_cdn_profile",
		F:    testSweepCDNProfiles,
	})
}

func testSweepCDNProfiles(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).cdnProfilesClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the CDN Profiles..")
	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on CDN Profiles: %+v", err)
	}

	for _, profile := range results.Values() {
		if !shouldSweepAcceptanceTestResource(*profile.Name, *profile.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*profile.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["profiles"]

		log.Printf("Deleting CDN Profile '%s' in Resource Group '%s'", name, resourceGroup)
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMCdnProfile_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMCdnProfile_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists("azurerm_cdn_profile.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnProfile_withTags(t *testing.T) {
	resourceName := "azurerm_cdn_profile.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMCdnProfile_withTags(ri, location)
	postConfig := testAccAzureRMCdnProfile_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnProfile_NonStandardCasing(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMCdnProfileNonStandardCasing(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists("azurerm_cdn_profile.test"),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMCdnProfile_basicToStandardAkamai(t *testing.T) {
	resourceName := "azurerm_cdn_profile.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMCdnProfile_basic(ri, testLocation())
	postConfig := testAccAzureRMCdnProfile_standardAkamai(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard_Verizon"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard_Akamai"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnProfile_standardAkamai(t *testing.T) {
	resourceName := "azurerm_cdn_profile.test"
	ri := acctest.RandInt()
	config := testAccAzureRMCdnProfile_standardAkamai(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Standard_Akamai"),
				),
			},
		},
	})
}

func testCheckAzureRMCdnProfileExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for cdn profile: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).cdnProfilesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cdnProfilesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CDN Profile %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMCdnProfileDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).cdnProfilesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cdn_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("CDN Profile still exists:\n%#v", resp.ProfileProperties)
		}
	}

	return nil
}

func testAccAzureRMCdnProfile_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"
}
`, rInt, location, rInt)
}

func testAccAzureRMCdnProfile_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCdnProfile_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Verizon"

  tags {
    environment = "staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCdnProfileNonStandardCasing(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "standard_verizon"
}
`, rInt, location, rInt)
}

func testAccAzureRMCdnProfile_standardAkamai(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard_Akamai"
}
`, rInt, location, rInt)
}
