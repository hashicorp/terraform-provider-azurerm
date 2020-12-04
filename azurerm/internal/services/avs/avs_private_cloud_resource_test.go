package avs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AvsPrivateCloudResource struct {
}

func TestAccAzureRMavsPrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	r := AvsPrivateCloudResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMavsPrivateCloud_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	r := AvsPrivateCloudResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMavsPrivateCloud_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	r := AvsPrivateCloudResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep("nsxt_password", "vcenter_password"),
	})
}

// Internet availability, cluster size, identity sources, vcenter password or nsxt password cannot be updated at the same time
func TestAccAzureRMavsPrivateCloud_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
	r := AvsPrivateCloudResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep("nsxt_password", "vcenter_password"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep("nsxt_password", "vcenter_password"),
		{
			Config: r.update2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep("nsxt_password", "vcenter_password"),
		{
			Config: r.update3(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
		data.ImportStep("nsxt_password", "vcenter_password"),
	})
}

func (AvsPrivateCloudResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PrivateCloudID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Avs.PrivateCloudClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Avs Private Cloud %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.PrivateCloudProperties != nil), nil
}

func (AvsPrivateCloudResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  # In Avs acctest, please disable correlation request id, else the continuous operations like update or delete will not be triggered
  disable_correlation_request_id = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-AVS-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AvsPrivateCloudResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-APC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    cluster_size = 3
  }
  network_block = "192.168.48.0/22"
}
`, r.template(data), data.RandomInteger)
}

func (r AvsPrivateCloudResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (r AvsPrivateCloudResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-APC-%d"
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
`, r.template(data), data.RandomInteger)
}

func (r AvsPrivateCloudResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-APC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "AV36"

  management_cluster {
    cluster_size = 4
  }
  network_block      = "192.168.48.0/22"
  internet_connected = false
  nsxt_password      = "QazWsx13$Edc"
  vcenter_password   = "QazWsx13$Edc"
  tags = {
    ENV = "Stage"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AvsPrivateCloudResource) update2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-APC-%d"
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
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r AvsPrivateCloudResource) update3(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-APC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "AV36"

  management_cluster {
    cluster_size = 3
  }
  network_block      = "192.168.48.0/22"
  internet_connected = true
  nsxt_password      = "QazWsx13$Edc"
  vcenter_password   = "QazWsx13$Edc"
  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}
