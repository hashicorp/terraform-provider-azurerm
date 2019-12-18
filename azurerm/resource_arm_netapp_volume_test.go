package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppVolume_basic(t *testing.T) {
	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
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

func TestAccAzureRMNetAppVolume_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetAppVolume_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_netapp_volume"),
			},
		},
	})
}

func TestAccAzureRMNetAppVolume_complete(t *testing.T) {
	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "storage_quota_in_gb", "101"),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "2"),
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

func TestAccAzureRMNetAppVolume_update(t *testing.T) {
	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_quota_in_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppVolume_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_quota_in_gb", "101"),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "storage_quota_in_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "0"),
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

func TestAccAzureRMNetAppVolume_updateSubnet(t *testing.T) {
	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resourceGroupName := fmt.Sprintf("acctestRG-netapp-%d", ri)
	oldVNetName := fmt.Sprintf("acctest-VirtualNetwork-%d", ri)
	oldSubnetName := fmt.Sprintf("acctest-Subnet-%d", ri)
	newVNetName := fmt.Sprintf("acctest-updated-VirtualNetwork-%d", ri)
	newSubnetName := fmt.Sprintf("acctest-updated-Subnet-%d", ri)
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
				Config: testAccAzureRMNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", oldSubnetId),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppVolume_updateSubnet(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", newSubnetId),
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

func TestAccAzureRMNetAppVolume_updateExportPolicyRule(t *testing.T) {
	resourceName := "azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetAppVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppVolume_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetAppVolume_updateExportPolicyRule(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppVolumeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "export_policy_rule.#", "1"),
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

func testCheckAzureRMNetAppVolumeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Volume not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).NetApp.VolumeClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMNetAppVolume_basic(rInt int, location string) string {
	template := testAccAzureRMNetAppVolume_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  storage_quota_in_gb = 100
}
`, template, rInt, rInt)
}

func testAccAzureRMNetAppVolume_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "import" {
  name                = "${azurerm_netapp_volume.test.name}"
  location            = "${azurerm_netapp_volume.test.location}"
  resource_group_name = "${azurerm_netapp_volume.test.name}"
}
`, testAccAzureRMNetAppVolume_basic(rInt, location))
}

func testAccAzureRMNetAppVolume_complete(rInt int, location string) string {
	template := testAccAzureRMNetAppVolume_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  service_level       = "Premium"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = "${azurerm_subnet.test.id}"
  storage_quota_in_gb = 101

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["1.2.3.0/24"]
    cifs_enabled    = false
    nfsv3_enabled   = true
    nfsv4_enabled   = false
    unix_read_only  = false
    unix_read_write = true
  }

  export_policy_rule {
    rule_index      = 2
    allowed_clients = ["1.2.5.0"]
    cifs_enabled    = false
    nfsv3_enabled   = true
    nfsv4_enabled   = false
    unix_read_only  = true
    unix_read_write = false
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMNetAppVolume_updateSubnet(rInt int, location string) string {
	template := testAccAzureRMNetAppVolume_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "updated" {
  name                = "acctest-updated-VirtualNetwork-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "updated" {
  name                 = "acctest-updated-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.updated.name}"
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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_path         = "my-updated-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.updated.id}"
  storage_quota_in_gb = 100
}
`, template, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetAppVolume_updateExportPolicyRule(rInt int, location string) string {
	template := testAccAzureRMNetAppVolume_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  service_level       = "Premium"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = "${azurerm_subnet.test.id}"
  storage_quota_in_gb = 101

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["1.2.4.0/24", "1.3.4.0"]
    cifs_enabled    = false
    nfsv3_enabled   = true
    nfsv4_enabled   = false
    unix_read_only  = false
    unix_read_write = true
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMNetAppVolume_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"

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
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_netapp_account.test.name}"
  service_level       = "Premium"
  size_in_tb          = 4
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
