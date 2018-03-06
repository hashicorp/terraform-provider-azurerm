package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMPublicIPs_namePrefix(t *testing.T) {
	dataSourceName := "data.azurerm_public_ips.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(5)

	config := testAccDataSourceAzureRMPublicIPs_prefix(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "public_ips.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpipa%s-0", rs)),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPs_assigned(t *testing.T) {
	attachedDataSourceName := "data.azurerm_public_ips.attached"
	unattachedDataSourceName := "data.azurerm_public_ips.unattached"
	ri := acctest.RandInt()
	rs := acctest.RandString(5)

	config := testAccDataSourceAzureRMPublicIPs_attached(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(attachedDataSourceName, "public_ips.#", "3"),
					resource.TestCheckResourceAttr(attachedDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpip%s-0", rs)),
					resource.TestCheckResourceAttr(unattachedDataSourceName, "public_ips.#", "4"),
					resource.TestCheckResourceAttr(unattachedDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpip%s-3", rs)),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIPs_allocationType(t *testing.T) {
	staticDataSourceName := "data.azurerm_public_ips.static"
	dynamicDataSourceName := "data.azurerm_public_ips.dynamic"
	ri := acctest.RandInt()
	rs := acctest.RandString(5)

	config := testAccDataSourceAzureRMPublicIPs_allocationType(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(staticDataSourceName, "public_ips.#", "3"),
					resource.TestCheckResourceAttr(staticDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpips%s-0", rs)),
					resource.TestCheckResourceAttr(dynamicDataSourceName, "public_ips.#", "4"),
					resource.TestCheckResourceAttr(dynamicDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpipd%s-0", rs)),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPublicIPs_attached(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                        = 7
  name                         = "acctestpip%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

resource "azurerm_lb" "test" {
  count               = 3
  name                = "acctestlb-${count.index}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = "${element(azurerm_public_ip.test.*.id, count.index)}"
  }
}

data "azurerm_public_ips" "unattached" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  attached            = false
  depends_on          = ["azurerm_lb.test"]
}

data "azurerm_public_ips" "attached" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  attached            = true
  depends_on          = ["azurerm_lb.test"]
}
`, rInt, location, rString)
}

func testAccDataSourceAzureRMPublicIPs_prefix(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                        = 2
  name                         = "acctestpipb%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

resource "azurerm_public_ip" "test2" {
  count                        = 2
  name                         = "acctestpipa%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

data "azurerm_public_ips" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name_prefix         = "acctestpipa"
  depends_on          = ["azurerm_public_ip.test", "azurerm_public_ip.test2"]
}
`, rInt, location, rString, rString)
}

func testAccDataSourceAzureRMPublicIPs_allocationType(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_public_ip" "dynamic" {
  count                        = 4
  name                         = "acctestpipd%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "dynamic"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

resource "azurerm_public_ip" "static" {
  count                        = 3
  name                         = "acctestpips%s-${count.index}"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

data "azurerm_public_ips" "dynamic" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_type     = "Dynamic"
  depends_on          = ["azurerm_public_ip.dynamic"]
}

data "azurerm_public_ips" "static" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_type     = "Static"
  depends_on          = ["azurerm_public_ip.static"]
}
`, rInt, location, rString, rString)
}
