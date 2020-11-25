package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
)

func TestAccAzureRMCdnProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMCdnProfile_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_cdn_profile"),
			},
		},
	})
}

func TestAccAzureRMCdnProfile_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCdnProfile_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnProfile_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfileNonStandardCasing(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists("azurerm_cdn_profile.test"),
				),
			},
			{
				Config:             testAccAzureRMCdnProfileNonStandardCasing(data),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAzureRMCdnProfile_basicToStandardAkamai(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard_Verizon"),
				),
			},
			{
				Config: testAccAzureRMCdnProfile_standardAkamai(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard_Akamai"),
				),
			},
		},
	})
}

func TestAccAzureRMCdnProfile_standardAkamai(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_standardAkamai(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard_Akamai"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCdnProfile_standardMicrosoft(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCdnProfile_standardMicrosoft(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard_Microsoft"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCdnProfileExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.ProfilesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ProfileID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cdnProfilesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CDN Profile %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMCdnProfileDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Cdn.ProfilesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cdn_profile" {
			continue
		}

		id, err := parse.ProfileID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("CDN Profile still exists:\n%#v", resp.ProfileProperties)
		}
	}

	return nil
}

func testAccAzureRMCdnProfile_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCdnProfile_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMCdnProfile_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_profile" "import" {
  name                = azurerm_cdn_profile.test.name
  location            = azurerm_cdn_profile.test.location
  resource_group_name = azurerm_cdn_profile.test.resource_group_name
  sku                 = azurerm_cdn_profile.test.sku
}
`, template)
}

func testAccAzureRMCdnProfile_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCdnProfile_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCdnProfileNonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standard_verizon"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCdnProfile_standardAkamai(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Akamai"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCdnProfile_standardMicrosoft(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
