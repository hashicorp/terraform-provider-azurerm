package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMPublicIPs_basic(t *testing.T) {
	dataSourceNameUsed := "data.azurerm_public_ips.test_used"
	dataSourceNameUnused := "data.azurerm_public_ips.test_unused"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPsBasic(name, resourceGroupName, testLocation(), 10, 0)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameUsed, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceNameUsed, "public_ips.#", "0"),
					resource.TestCheckResourceAttr(dataSourceNameUnused, "public_ips.#", "10"),
					resource.TestCheckResourceAttr(dataSourceNameUnused, "public_ips.0.name", fmt.Sprintf("%s-0", name)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPs_mixed(t *testing.T) {
	dataSourceNameUsed := "data.azurerm_public_ips.test_used"
	dataSourceNameUnused := "data.azurerm_public_ips.test_unused"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPsBasic(name, resourceGroupName, testLocation(), 10, 6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameUsed, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceNameUsed, "public_ips.#", "6"),
					resource.TestCheckResourceAttr(dataSourceNameUsed, "public_ips.0.name", fmt.Sprintf("%s-0", name)),
					resource.TestCheckResourceAttr(dataSourceNameUnused, "public_ips.#", "4"),
					resource.TestCheckResourceAttr(dataSourceNameUnused, "public_ips.0.name", fmt.Sprintf("%s-6", name)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPs_count(t *testing.T) {
	dataSourceNameUnused := "data.azurerm_public_ips.test_unused"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPsCount(name, resourceGroupName, testLocation(), 10, 5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameUnused, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceNameUnused, "public_ips.#", "10"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPs_tooFew(t *testing.T) {
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPsCount(name, resourceGroupName, testLocation(), 10, 15)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Not enough unassigned public IP addresses in resource group %s", resourceGroupName)),
			},
		},
	})
}

func randNames() (string, string) {
	ri := acctest.RandInt()
	name := fmt.Sprintf("acctestpublicippublic_ips-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)
	return name, resourceGroupName
}

func testAccDataSourceAzureRMPublicIPsBasic(name string, resourceGroupName string, location string, pipCount int, lbCount int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                        = %d
  name                         = "%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

resource "azurerm_lb" "test" {
  count               = %d
  name                = "load-balancer-${count.index}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "${element(azurerm_public_ip.test.*.id, count.index)}"
  }
}

data "azurerm_public_ips" "test_unused" {
	resource_group_name = "${azurerm_resource_group.test.name}"
	attached            = false
  depends_on          = ["azurerm_lb.test", "azurerm_public_ip.test"]
}
data "azurerm_public_ips" "test_used" {
	resource_group_name = "${azurerm_resource_group.test.name}"
	attached            = true
  depends_on          = ["azurerm_lb.test", "azurerm_public_ip.test"]
}
`, resourceGroupName, location, pipCount, name, lbCount)
}

func testAccDataSourceAzureRMPublicIPsCount(name string, resourceGroupName string, location string, pipCount int, minCount int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                        = %d
  name                         = "%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

data "azurerm_public_ips" "test_unused" {
	resource_group_name = "${azurerm_resource_group.test.name}"
	attached            = false
  minimum_count       = %d
  depends_on          = ["azurerm_public_ip.test"]
}
`, resourceGroupName, location, pipCount, name, minCount)
}
