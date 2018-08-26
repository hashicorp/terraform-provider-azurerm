package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMImage_basic(t *testing.T) {
	dataSourceName := "data.azurerm_image.test"

	config := testAccDataSourceAzureRMImageBasic(acctest.RandInt(), acctest.RandString(4), testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "os_disk.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "os_disk.0.blob_uri"),
					resource.TestCheckResourceAttr(dataSourceName, "os_disk.0.caching", "None"),
					resource.TestCheckResourceAttr(dataSourceName, "os_disk.0.os_type", "Linux"),
					resource.TestCheckResourceAttr(dataSourceName, "os_disk.0.os_state", "Generalized"),
					resource.TestCheckResourceAttr(dataSourceName, "os_disk.0.size_gb", "30"),
					resource.TestCheckResourceAttr(dataSourceName, "data_disk.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "Dev"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.cost-center", "Ops"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMImage_localFilter(t *testing.T) {
	ascDataSourceName := "data.azurerm_image.test1"
	descDataSourceName := "data.azurerm_image.test2"

	ri := acctest.RandInt()
	config := testAccDataSourceAzureRMImageLocalFilter(ri, acctest.RandString(4), testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ascDataSourceName, "name"),
					resource.TestCheckResourceAttrSet(ascDataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(ascDataSourceName, "name", fmt.Sprintf("def-acctest-%d", ri)),

					resource.TestCheckResourceAttrSet(descDataSourceName, "name"),
					resource.TestCheckResourceAttrSet(descDataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(descDataSourceName, "name", fmt.Sprintf("def-acctest-%d", ri)),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMImageBasic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "acctestpip%d"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctestnic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "acctestvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.testsource.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "acctest-%d"
    admin_username = "tfuser"
    admin_password = "P@ssW0RD7890"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_virtual_machine.testsource.storage_os_disk.0.vhd_uri}"
    size_gb  = 30
    caching  = "None"
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_image" "test" {
  name                = "${azurerm_image.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

output "location" {
  value = "${data.azurerm_image.test.location}"
}

`, rInt, location, rInt, rInt, rInt, rInt, rString, rInt, rInt, rInt)
}

func testAccDataSourceAzureRMImageLocalFilter(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "acctestpip%d"
}

resource "azurerm_network_interface" "testsource" {
  name                = "acctestnic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfigurationsource"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "Dev"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "blob"
}

resource "azurerm_virtual_machine" "testsource" {
  name                  = "acctestvm-%d"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.testsource.id}"]
  vm_size               = "Standard_D1_v2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = "30"
  }

  os_profile {
    computer_name  = "acctest-%d"
    admin_username = "tfuser"
    admin_password = "P@ssW0RD7890"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

resource "azurerm_image" "abc" {
  name                = "abc-acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_virtual_machine.testsource.storage_os_disk.0.vhd_uri}"
    size_gb  = 30
    caching  = "None"
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}


resource "azurerm_image" "def" {
  name                = "def-acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = "${azurerm_virtual_machine.testsource.storage_os_disk.0.vhd_uri}"
    size_gb  = 30
    caching  = "None"
  }

  tags {
    environment = "Dev"
    cost-center = "Ops"
  }
}

data "azurerm_image" "test1" {
  name_regex          = "^def-acctest-\\d+"
  resource_group_name = "${azurerm_resource_group.test.name}${substr(azurerm_image.abc.name, 0, 0)}${substr(azurerm_image.def.name, 0, 0)}"
}

data "azurerm_image" "test2" {
  name_regex          = "^[a-z]+-acctest-\\d+"
  sort_descending     = true
  resource_group_name = "${azurerm_resource_group.test.name}${substr(azurerm_image.abc.name, 0, 0)}${substr(azurerm_image.def.name, 0, 0)}"
}

output "location" {
  value = "${data.azurerm_image.test1.location}"
}

`, rInt, location, rInt, rInt, rInt, rInt, rString, rInt, rInt, rInt, rInt)
}
