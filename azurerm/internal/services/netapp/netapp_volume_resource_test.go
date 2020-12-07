package netapp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.2676449260", "NFSv3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppVolume_nfsv41(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_nfsv41(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.3098200649", "NFSv4.1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppVolume_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppVolume_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_volume"),
			},
		},
	})
}

func TestAccAzureRMNetAppVolume_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_level", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_quota_in_gb", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
					resource.TestCheckResourceAttr(data.ResourceName, "mount_ip_addresses.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppVolume_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_quota_in_gb", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppVolume_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_quota_in_gb", "101"),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.FoO", "BaR"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_quota_in_gb", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppVolume_updateSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	resourceGroupName := fmt.Sprintf("acctestRG-netapp-%d", data.RandomInteger)
	oldVNetName := fmt.Sprintf("acctest-VirtualNetwork-%d", data.RandomInteger)
	oldSubnetName := fmt.Sprintf("acctest-Subnet-%d", data.RandomInteger)
	newVNetName := fmt.Sprintf("acctest-updated-VirtualNetwork-%d", data.RandomInteger)
	newSubnetName := fmt.Sprintf("acctest-updated-Subnet-%d", data.RandomInteger)
	uriTemplate := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s"

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	oldSubnetId := fmt.Sprintf(uriTemplate, subscriptionID, resourceGroupName, oldVNetName, oldSubnetName)
	newSubnetId := fmt.Sprintf(uriTemplate, subscriptionID, resourceGroupName, newVNetName, newSubnetName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_id", oldSubnetId),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppVolume_updateSubnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "subnet_id", newSubnetId),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetAppVolume_updateExportPolicyRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "3"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetAppVolume_updateExportPolicyRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "export_policy_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetAppVolumeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.VolumeClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Volume not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Volume %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.VolumeClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppVolumeDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.VolumeClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_volume" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.VolumeClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppVolume_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetAppVolume_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppVolume_nfsv41(data acceptance.TestData) string {
	template := testAccAzureRMNetAppVolume_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv4.1"]
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv4.1"]
    unix_read_only    = false
    unix_read_write   = true
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppVolume_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "import" {
  name                = azurerm_netapp_volume.test.name
  location            = azurerm_netapp_volume.test.location
  resource_group_name = azurerm_netapp_volume.test.resource_group_name
  account_name        = azurerm_netapp_volume.test.account_name
  pool_name           = azurerm_netapp_volume.test.pool_name
  volume_path         = azurerm_netapp_volume.test.volume_path
  service_level       = azurerm_netapp_volume.test.service_level
  subnet_id           = azurerm_netapp_volume.test.subnet_id
  storage_quota_in_gb = azurerm_netapp_volume.test.storage_quota_in_gb
}
`, testAccAzureRMNetAppVolume_basic(data))
}

func testAccAzureRMNetAppVolume_complete(data acceptance.TestData) string {
	template := testAccAzureRMNetAppVolume_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 101

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  export_policy_rule {
    rule_index        = 2
    allowed_clients   = ["1.2.5.0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  export_policy_rule {
    rule_index      = 3
    allowed_clients = ["1.2.6.0/24"]
    cifs_enabled    = false
    nfsv3_enabled   = true
    nfsv4_enabled   = false
    unix_read_only  = true
    unix_read_write = false
  }

  tags = {
    "FoO" = "BaR"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppVolume_updateSubnet(data acceptance.TestData) string {
	template := testAccAzureRMNetAppVolume_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "updated" {
  name                = "acctest-updated-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "updated" {
  name                 = "acctest-updated-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.updated.name
  address_prefix       = "10.1.3.0/24"

  delegation {
    name = "testdelegation2"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-updated-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-updated-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.updated.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 100
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppVolume_updateExportPolicyRule(data acceptance.TestData) string {
	template := testAccAzureRMNetAppVolume_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 101

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.4.0/24", "1.3.4.0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "FoO" = "BaR"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetAppVolume_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.6.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.6.2.0/24"

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
