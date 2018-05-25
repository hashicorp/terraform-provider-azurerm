package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceAzureRMPublicIpDomainNameLabel_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting123",
			ErrCount: 1,
		},
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandString(80),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validatePublicIpDomainNameLabel(tc.Value, "azurerm_public_ip")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Public IP Domain Name Label to trigger a validation error")
		}
	}
}

func TestAccAzureRMPublicIpStatic_basic(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPublicIPStatic_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address_allocation", "static"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_basic_withDNSLabel(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := acctest.RandInt()
	dnl := fmt.Sprintf("acctestdnl-%d", ri)
	config := testAccAzureRMPublicIPStatic_basic_withDNSLabel(ri, testLocation(), dnl)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address_allocation", "static"),
					resource.TestCheckResourceAttr(resourceName, "domain_name_label", dnl),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_standard(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMPublicIPStatic_standard(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists("azurerm_public_ip.test"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_disappears(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := acctest.RandInt()
	config := testAccAzureRMPublicIPStatic_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	config := testAccAzureRMPublicIPStatic_idleTimeout(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout_in_minutes", "30"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpStatic_withTags(t *testing.T) {
	resourceName := "azurerm_public_ip.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMPublicIPStatic_withTags(ri, location)
	postConfig := testAccAzureRMPublicIPStatic_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMPublicIPStatic_basic(ri, location)
	postConfig := testAccAzureRMPublicIPStatic_update(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		},
	})
}

func TestAccAzureRMPublicIpDynamic_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMPublicIPDynamic_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIpExists("azurerm_public_ip.test"),
				),
			},
		},
	})
}

func testCheckAzureRMPublicIpExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		publicIPName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip: %s", publicIPName)
		}

		client := testAccProvider.Meta().(*ArmClient).publicIPClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, publicIPName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIPClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Public IP %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPublicIpDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		publicIpName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip: %s", publicIpName)
		}

		client := testAccProvider.Meta().(*ArmClient).publicIPClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, publicIpName)
		if err != nil {
			return fmt.Errorf("Error deleting Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for deletion of Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMPublicIpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).publicIPClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

    public_ip_address_allocation = "static"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_basic_withZone(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name                         = "acctestpublicip-%d"
    location                     = "${azurerm_resource_group.test.location}"
    resource_group_name          = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"

    zones = ["1"]
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_basic_withDNSLabel(rInt int, location, dnsNameLabel string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name        = "acctestRG-%d"
    location    = "%s"
}

resource "azurerm_public_ip" "test" {
    name                            = "acctestpublicip-%d"
    location                        = "${azurerm_resource_group.test.location}"
    resource_group_name             = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation    = "static"

    domain_name_label = "%s"
}
`, rInt, location, rInt, dnsNameLabel)
}

func testAccAzureRMPublicIPStatic_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
    sku = "standard"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
    domain_name_label = "acctest-%d"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPublicIPStatic_idleTimeout(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
    idle_timeout_in_minutes = 30
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPDynamic_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "dynamic"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPStatic_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_public_ip" "test" {
    name = "acctestpublicip-%d"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"

    tags {
	environment = "staging"
    }
}
`, rInt, location, rInt)
}
