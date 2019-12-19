package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMPublicIpStatic_basic(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
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

func TestAccAzureRMPublicIpStatic_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
				),
			},
			{
				Config:      testAccAzureRMPublicIPStatic_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_public_ip"),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_basicOld(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPStatic_basicOld(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address_allocation", "Static"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
				),
			},
			{
				Config: testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
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

func TestAccAzureRMPublicIpStatic_zones(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPStatic_withZone(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "zones.0", "1"),
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

func TestAccAzureRMPublicIpStatic_basic_withDNSLabel(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	dnl := fmt.Sprintf("acctestdnl-%d", ri)
	config := testAccAzureRMPublicIPStatic_basic_withDNSLabel(ri, acceptance.Location(), dnl)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_label", dnl),
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

func TestAccAzureRMPublicIpStatic_standard_withIPv6_fails(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_standard_withIPVersion(ri, acceptance.Location(), "IPv6")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Cannot specify publicIpAllocationMethod as Static for IPv6 PublicIp"),
			},
		},
	})
}

func TestAccAzureRMPublicIpDynamic_basic_withIPv6(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	ipVersion := "Ipv6"
	config := testAccAzureRMPublicIPDynamic_basic_withIPVersion(ri, acceptance.Location(), ipVersion)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv6"),
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

func TestAccAzureRMPublicIpStatic_basic_defaultsToIPv4(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
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
func TestAccAzureRMPublicIpStatic_basic_withIPv4(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	ipVersion := "IPv4"
	config := testAccAzureRMPublicIPStatic_basic_withIPVersion(ri, acceptance.Location(), ipVersion)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "IPv4"),
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

func TestAccAzureRMPublicIpStatic_standard(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_standard(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
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

func TestAccAzureRMPublicIpStatic_disappears(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					testCheckAzureRMPublicIpDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_idleTimeout(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_idleTimeout(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout_in_minutes", "30"),
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

func TestAccAzureRMPublicIpStatic_withTags(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMPublicIPStatic_withTags(ri, location)
	postConfig := testAccAzureRMPublicIPStatic_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_update(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMPublicIPStatic_basic(ri, location)
	postConfig := testAccAzureRMPublicIPStatic_update(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "domain_name_label", fmt.Sprintf("acctest-%d", ri)),
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

func TestAccAzureRMPublicIpStatic_standardPrefix(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_standardPrefix(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_standardPrefixWithTags(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMPublicIPStatic_standardPrefixWithTags(ri, location)
	postConfig := testAccAzureRMPublicIPStatic_standardPrefixWithTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpDynamic_basic(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPDynamic_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
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

func TestAccAzureRMPublicIpStatic_importIdError(t *testing.T) {
	resourceName := "azurerm_public_ip.test"

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPStatic_basic(ri, acceptance.Location())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/publicIPAdresses/acctestpublicip-%d", os.Getenv("ARM_SUBSCRIPTION_ID"), ri, ri),
				ExpectError:       regexp.MustCompile("Error parsing supplied resource id."),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_canLabelBe63(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPStatic_canLabelBe63(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "allocation_method", "Static"),
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

func testCheckAzureRMPublicIpExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIPName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip: %s", publicIPName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, publicIPName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIPClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Public IP %q (resource group: %q) does not exist", publicIPName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPublicIpDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip: %s", publicIpName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		future, err := client.Delete(ctx, resourceGroup, publicIpName)
		if err != nil {
			return fmt.Errorf("Error deleting Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMPublicIpDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_public_ip" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Public IP still exists:\n%#v", resp.PublicIPAddressPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMPublicIPStatic_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "import" {
  name                = "${azurerm_public_ip.test.name}"
  location            = "${azurerm_public_ip.test.location}"
  resource_group_name = "${azurerm_public_ip.test.resource_group_name}"
  allocation_method   = "${azurerm_public_ip.test.allocation_method}"
}
`, testAccAzureRMPublicIPStatic_basic(rInt, location))
}

func testAccAzureRMPublicIPStatic_basicOld(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpublicip-%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_withZone(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  zones               = ["1"]
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_basic_withDNSLabel(rInt int, location, dnsNameLabel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = "%s"
}
`, rInt, location, rInt, dnsNameLabel)
}

func testAccAzureRMPublicIPStatic_basic_withIPVersion(rInt int, location string, ipVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  ip_version          = "%s"
}
`, rInt, location, rInt, ipVersion)
}

func testAccAzureRMPublicIPStatic_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_standardPrefix(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = "${azurerm_public_ip_prefix.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPublicIPStatic_standardPrefixWithTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = "${azurerm_public_ip_prefix.test.id}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPublicIPStatic_standardPrefixWithTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = "${azurerm_public_ip_prefix.test.id}"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPublicIPStatic_standard_withIPVersion(rInt int, location string, ipVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  ip_version          = "%s"
  sku                 = "Standard"
}
`, rInt, location, rInt, ipVersion)
}

func testAccAzureRMPublicIPStatic_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  domain_name_label   = "acctest-%d"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPublicIPStatic_idleTimeout(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPDynamic_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPDynamic_basic_withIPVersion(rInt int, location string, ipVersion string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Dynamic"

  ip_version = "%s"
}
`, rInt, location, rInt, ipVersion)
}

func testAccAzureRMPublicIPStatic_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_canLabelBe63(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  allocation_method = "Static"
  domain_name_label = "k2345678-1-2345678-2-2345678-3-2345678-4-2345678-5-2345678-6-23"
}
`, rInt, location, rInt)
}
