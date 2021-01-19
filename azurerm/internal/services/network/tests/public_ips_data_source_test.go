package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PublicIPsResource struct {
}

func TestAccDataSourcePublicIPs_namePrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ips", "test")
	r := PublicIPsResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.prefix(data),
		},
		{
			Config: r.prefixDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_ips.#").HasValue("2"),
				check.That(data.ResourceName).Key("public_ips.0.name").HasValue(fmt.Sprintf("acctestpipa%s-0", data.RandomString)),
			),
		},
	})
}

func TestAccDataSourcePublicIPs_assigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ips", "test")
	r := PublicIPsResource{}

	attachedDataSourceName := "data.azurerm_public_ips.attached"
	unattachedDataSourceName := "data.azurerm_public_ips.unattached"

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.attached(data),
		},
		{
			Config: r.attachedDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(attachedDataSourceName, "public_ips.#", "3"),
				resource.TestCheckResourceAttr(attachedDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpip%s-0", data.RandomString)),
				resource.TestCheckResourceAttr(unattachedDataSourceName, "public_ips.#", "4"),
				resource.TestCheckResourceAttr(unattachedDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpip%s-3", data.RandomString)),
			),
		},
	})
}

func TestAccDataSourcePublicIPs_allocationType(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ips", "test")
	r := PublicIPsResource{}

	staticDataSourceName := "data.azurerm_public_ips.static"
	dynamicDataSourceName := "data.azurerm_public_ips.dynamic"

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.allocationType(data),
		},
		{
			Config: r.allocationTypeDataSources(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(staticDataSourceName, "public_ips.#", "3"),
				resource.TestCheckResourceAttr(staticDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpips%s-0", data.RandomString)),
				resource.TestCheckResourceAttr(dynamicDataSourceName, "public_ips.#", "4"),
				resource.TestCheckResourceAttr(dynamicDataSourceName, "public_ips.0.name", fmt.Sprintf("acctestpipd%s-0", data.RandomString)),
			),
		},
	})
}

func (PublicIPsResource) attached(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                   = 7
  name                    = "acctestpip%s-${count.index}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_lb" "test" {
  count               = 3
  name                = "acctestlb-${count.index}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "frontend"
    public_ip_address_id = element(azurerm_public_ip.test.*.id, count.index)
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r PublicIPsResource) attachedDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "unattached" {
  resource_group_name = azurerm_resource_group.test.name
  attached            = false
}

data "azurerm_public_ips" "attached" {
  resource_group_name = azurerm_resource_group.test.name
  attached            = true
}
`, r.attached(data))
}

func (PublicIPsResource) prefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  count                   = 2
  name                    = "acctestpipb%s-${count.index}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_public_ip" "test2" {
  count                   = 2
  name                    = "acctestpipa%s-${count.index}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r PublicIPsResource) prefixDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name_prefix         = "acctestpipa"
}
`, r.prefix(data))
}

func (PublicIPsResource) allocationType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "dynamic" {
  count                   = 4
  name                    = "acctestpipd%s-${count.index}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Dynamic"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

resource "azurerm_public_ip" "static" {
  count                   = 3
  name                    = "acctestpips%s-${count.index}"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r PublicIPsResource) allocationTypeDataSources(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_public_ips" "dynamic" {
  resource_group_name = azurerm_resource_group.test.name
  allocation_type     = "Dynamic"
}

data "azurerm_public_ips" "static" {
  resource_group_name = azurerm_resource_group.test.name
  allocation_type     = "Static"
}
`, r.allocationType(data))
}
