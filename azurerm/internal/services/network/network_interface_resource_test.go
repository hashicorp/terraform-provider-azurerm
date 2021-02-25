package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkInterfaceResource struct {
}

func TestAccNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
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

func TestAccNetworkInterface_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccNetworkInterface_dnsServers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dnsServers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dnsServersUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_enableAcceleratedNetworking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Enabled
			Config: r.enableAcceleratedNetworking(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.enableAcceleratedNetworking(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.enableAcceleratedNetworking(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_enableIPForwarding(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Enabled
			Config: r.enableIPForwarding(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.enableIPForwarding(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.enableIPForwarding(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_internalDomainNameLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.internalDomainNameLabel(data, "1"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.internalDomainNameLabel(data, "2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ipv6(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("ip_configuration.1.private_ip_address_version").HasValue("IPv6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_multipleIPConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleIPConfigurations(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_multipleIPConfigurationsSecondaryAsPrimary(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleIPConfigurationsSecondaryAsPrimary(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_publicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.publicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicIPRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.publicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_interface"),
		},
	})
}

func TestAccNetworkInterface_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleIPConfigurations(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkInterface_updateMultipleParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface", "test")
	r := NetworkInterfaceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withMultipleParameters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateMultipleParameters(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkInterfaceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NetworkInterfaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.InterfacesClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (NetworkInterfaceResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NetworkInterfaceID(state.ID)
	if err != nil {
		return nil, err
	}

	future, err := client.Network.InterfacesClient.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Network.InterfacesClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r NetworkInterfaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) withMultipleParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r NetworkInterfaceResource) updateMultipleParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  enable_ip_forwarding    = true
  internal_dns_name_label = "acctestni-%s"

  dns_servers = [
    "10.0.0.5",
    "10.0.0.7"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r NetworkInterfaceResource) dnsServers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.5",
    "10.0.0.6"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) dnsServersUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  dns_servers = [
    "10.0.0.6",
    "10.0.0.5"
  ]

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) enableAcceleratedNetworking(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                          = "acctestni-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  enable_accelerated_networking = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r NetworkInterfaceResource) enableIPForwarding(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                 = "acctestni-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  enable_ip_forwarding = %t

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, enabled)
}

func (r NetworkInterfaceResource) internalDomainNameLabel(data acceptance.TestData, suffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                    = "acctestni-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  internal_dns_name_label = "acctestni-%s-%s"

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger, suffix, data.RandomString)
}

func (r NetworkInterfaceResource) ipv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    private_ip_address_allocation = "Dynamic"
    private_ip_address_version    = "IPv6"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) multipleIPConfigurations(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) multipleIPConfigurationsSecondaryAsPrimary(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  ip_configuration {
    name                          = "secondary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}
`, r.publicIPTemplate(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIPRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.publicIPTemplate(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) publicIPTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "import" {
  name                = azurerm_network_interface.test.name
  location            = azurerm_network_interface.test.location
  resource_group_name = azurerm_network_interface.test.resource_group_name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, r.basic(data))
}

func (r NetworkInterfaceResource) static(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.2.15"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkInterfaceResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "primary"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Hello     = "World"
    Elephants = "Five"
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetworkInterfaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
