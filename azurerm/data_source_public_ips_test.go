package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPublicIPs_namePrefix(t *testing.T) {
	dataSourceName := "data.azurerm_public_ips.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resourceConfig := testAccDataSourceAzureRMPublicIPs_prefix(ri, rs, location)
	dataSourceConfig := testAccDataSourceAzureRMPublicIPs_prefixDataSource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: dataSourceConfig,
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
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resourceConfig := testAccDataSourceAzureRMPublicIPs_attached(ri, rs, location)
	dataSourceConfig := testAccDataSourceAzureRMPublicIPs_attachedDataSource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: dataSourceConfig,
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
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resourceConfig := testAccDataSourceAzureRMPublicIPs_allocationType(ri, rs, location)
	dataSourceConfig := testAccDataSourceAzureRMPublicIPs_allocationTypeDataSources(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: dataSourceConfig,
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                   = 7
  name                    = "acctestpip%s-${count.index}"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
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
`, rInt, location, rString)
}

func testAccDataSourceAzureRMPublicIPs_attachedDataSource(rInt int, rString string, location string) string {
	resources := testAccDataSourceAzureRMPublicIPs_attached(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "unattached" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  attached            = false
}

data "azurerm_public_ips" "attached" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  attached            = true
}
`, resources)
}

func testAccDataSourceAzureRMPublicIPs_prefix(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                   = 2
  name                    = "acctestpipb%s-${count.index}"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_public_ip" "test2" {
  count                   = 2
  name                    = "acctestpipa%s-${count.index}"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}
`, rInt, location, rString, rString)
}

func testAccDataSourceAzureRMPublicIPs_prefixDataSource(rInt int, rString string, location string) string {
	prefixed := testAccDataSourceAzureRMPublicIPs_prefix(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name_prefix         = "acctestpipa"
}
`, prefixed)
}

func testAccDataSourceAzureRMPublicIPs_allocationType(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "dynamic" {
  count                   = 4
  name                    = "acctestpipd%s-${count.index}"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_public_ip" "static" {
  count                   = 3
  name                    = "acctestpips%s-${count.index}"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}
`, rInt, location, rString, rString)
}

func testAccDataSourceAzureRMPublicIPs_allocationTypeDataSources(rInt int, rString string, location string) string {
	allocationType := testAccDataSourceAzureRMPublicIPs_allocationType(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "dynamic" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_type     = "Dynamic"
}

data "azurerm_public_ips" "static" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_type     = "Static"
}
`, allocationType)
}
