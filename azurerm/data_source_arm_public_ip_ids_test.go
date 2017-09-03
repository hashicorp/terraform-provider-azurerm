package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMPublicIPIds_basic(t *testing.T) {
	dataSourceName := "data.azurerm_public_ip_ids.test"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPIdsBasic(name, resourceGroupName, testLocation(), 10, 0)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "10"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPIds_mixed(t *testing.T) {
	dataSourceName := "data.azurerm_public_ip_ids.test"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPIdsBasic(name, resourceGroupName, testLocation(), 10, 6)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "4"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPIds_count(t *testing.T) {
	dataSourceName := "data.azurerm_public_ip_ids.test"
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPIdsCount(name, resourceGroupName, testLocation(), 10, 5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "10"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPIds_tooFew(t *testing.T) {
	name, resourceGroupName := randNames()

	config := testAccDataSourceAzureRMPublicIPIdsCount(name, resourceGroupName, testLocation(), 10, 15)

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
	name := fmt.Sprintf("acctestpublicipids-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)
	return name, resourceGroupName
}

func testAccDataSourceAzureRMPublicIPIdsBasic(name string, resourceGroupName string, location string, pipCount int, lbCount int) string {
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

data "azurerm_public_ip_ids" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  depends_on          = ["azurerm_lb.test", "azurerm_public_ip.test"]
}
`, resourceGroupName, location, pipCount, name, lbCount)
}

func testAccDataSourceAzureRMPublicIPIdsCount(name string, resourceGroupName string, location string, pipCount int, minCount int) string {
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

data "azurerm_public_ip_ids" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  minimum_count       = %d
  depends_on          = ["azurerm_public_ip.test"]
}
`, resourceGroupName, location, pipCount, name, minCount)
}
