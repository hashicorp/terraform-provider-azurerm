package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMavsPrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsPrivateCloud_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.hosts.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_private_peering_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.primary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.secondary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hcx_cloud_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "provisioning_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcenter_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcsa_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vmotion_network"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMavsPrivateCloud_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsPrivateCloud_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMavsPrivateCloud_requiresImport),
		},
	})
}

func TestAccAzureRMavsPrivateCloud_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsPrivateCloud_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.hosts.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_private_peering_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.primary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.secondary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hcx_cloud_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "provisioning_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcenter_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcsa_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vmotion_network"),
				),
			},
			data.ImportStep("nsxt_password", "vcenter_password"),
		},
	})
}

func TestAccAzureRMavsPrivateCloud_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsPrivateCloud_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.hosts.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_private_peering_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.primary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.secondary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hcx_cloud_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "provisioning_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcenter_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcsa_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vmotion_network"),
				),
			},
			data.ImportStep("nsxt_password", "vcenter_password"),
			{
				Config: testAccAzureRMavsPrivateCloud_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.hosts.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_private_peering_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.primary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.secondary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hcx_cloud_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "provisioning_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcenter_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcsa_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vmotion_network"),
				),
			},
			data.ImportStep("nsxt_password", "vcenter_password"),
		},
	})
}

func testCheckAzureRMavsPrivateCloudExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.PrivateCloudClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("avs PrivateCloud not found: %s", resourceName)
		}
		id, err := parse.AvsPrivateCloudID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Avs PrivateCloud %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Avs.PrivateCloudClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMavsPrivateCloudDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.PrivateCloudClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_avs_private_cloud" {
			continue
		}
		id, err := parse.AvsPrivateCloudID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Avs.PrivateCloudClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMavsPrivateCloud_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  # In Avs, we can't use client with the same x-ms-correlation-request-id for delete, else the delete will not be triggered
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-avs-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMavsPrivateCloud_basic(data acceptance.TestData) string {
	template := testAccAzureRMavsPrivateCloud_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    cluster_size = 3
  }
  network_block = "192.168.48.0/22"
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsPrivateCloud_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMavsPrivateCloud_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "import" {
  name                = azurerm_avs_private_cloud.test.name
  resource_group_name = azurerm_avs_private_cloud.test.resource_group_name
  location            = azurerm_avs_private_cloud.test.location
  sku_name            = azurerm_avs_private_cloud.test.sku_name

  management_cluster {
    cluster_size = azurerm_avs_private_cloud.test.management_cluster.0.cluster_size
  }
  network_block = azurerm_avs_private_cloud.test.network_block
}
`, config)
}

func testAccAzureRMavsPrivateCloud_complete(data acceptance.TestData) string {
	template := testAccAzureRMavsPrivateCloud_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "AV36"

  management_cluster {
    cluster_size = 3
  }
  network_block      = "192.168.48.0/22"
  internet_connected = false
  nsxt_password      = "QazWsx13$Edc"
  vcenter_password   = "QazWsx13$Edc"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsPrivateCloud_update(data acceptance.TestData) string {
	template := testAccAzureRMavsPrivateCloud_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "AV36"

  management_cluster {
    cluster_size = 4
  }
  network_block      = "192.168.48.0/22"
  internet_connected = true
  nsxt_password      = "QazWsx13$Edc"
  vcenter_password   = "QazWsx13$Edc"
  tags = {
    ENV = "Stage"
  }
}
`, template, data.RandomInteger)
}
