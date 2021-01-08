package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KustoClusterResource struct {
}

func TestAccKustoCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_disk_encryption").HasValue("false"),
				check.That(data.ResourceName).Key("enable_streaming_ingest").HasValue("false"),
				check.That(data.ResourceName).Key("enable_purge").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_disk_encryption").HasValue("true"),
				check.That(data.ResourceName).Key("enable_streaming_ingest").HasValue("true"),
				check.That(data.ResourceName).Key("enable_purge").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_disk_encryption").HasValue("false"),
				check.That(data.ResourceName).Key("enable_streaming_ingest").HasValue("false"),
				check.That(data.ResourceName).Key("enable_purge").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
	})
}

func TestAccKustoCluster_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Dev(No SLA)_Standard_D11_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("1"),
			),
		},
		{
			Config: r.skuUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku.0.name").HasValue("Standard_D11_v2"),
				check.That(data.ResourceName).Key("sku.0.capacity").HasValue("2"),
			),
		},
	})
}

func TestAccKustoCluster_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withZones(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
			),
		},
	})
}

func TestAccKustoCluster_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.engine_public_ip_id").Exists(),
				check.That(data.ResourceName).Key("virtual_network_configuration.0.data_management_public_ip_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_languageExtensions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.languageExtensions(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("language_extensions.#").HasValue("2"),
				check.That(data.ResourceName).Key("language_extensions.0").HasValue("PYTHON"),
				check.That(data.ResourceName).Key("language_extensions.1").HasValue("R"),
			),
		},
		data.ImportStep(),
		{
			Config: r.languageExtensionsRemove(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("language_extensions.#").HasValue("1"),
				check.That(data.ResourceName).Key("language_extensions.0").HasValue("R"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_optimizedAutoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.optimizedAutoScale(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("optimized_auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.minimum_instances").HasValue("2"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.maximum_instances").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.optimizedAutoScaleUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("optimized_auto_scale.#").HasValue("1"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.minimum_instances").HasValue("3"),
				check.That(data.ResourceName).Key("optimized_auto_scale.0.maximum_instances").HasValue("4"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCluster_engineV3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cluster", "test")
	r := KustoClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.engineV3(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KustoClusterResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.ClustersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.ClusterProperties != nil), nil
}

func (KustoClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) skuUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_D11_v2"
    capacity = 2
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) withZones(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  zones = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                    = "acctestkc%s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_disk_encryption  = true
  enable_streaming_ingest = true
  enable_purge            = true

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) languageExtensions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  language_extensions = ["PYTHON", "R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) languageExtensionsRemove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  language_extensions = ["R"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) optimizedAutoScale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_D11_v2"
  }

  optimized_auto_scale {
    minimum_instances = 2
    maximum_instances = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) optimizedAutoScaleUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name = "Standard_D11_v2"
  }

  optimized_auto_scale {
    minimum_instances = 3
    maximum_instances = 4
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KustoClusterResource) vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestkc%s-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestkc%s-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestkc%s-nsg"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_rule" "test_allow_management_inbound" {
  name                        = "AllowAzureDataExplorerManagement"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "AzureDataExplorerManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_public_ip" "engine_pip" {
  name                = "acctestkc%s-engine-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  allocation_method   = "Static"
}

resource "azurerm_public_ip" "management_pip" {
  name                = "acctestkc%s-management-pip"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  allocation_method   = "Static"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  virtual_network_configuration {
    subnet_id                    = azurerm_subnet.test.id
    engine_public_ip_id          = azurerm_public_ip.engine_pip.id
    data_management_public_ip_id = azurerm_public_ip.management_pip.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (KustoClusterResource) engineV3(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
  engine = "V3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
