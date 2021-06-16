package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PrivateLinkServiceResource struct {
}

func TestAccPrivateLinkService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestPLS-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateLinkService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_private_link_service"),
		},
	})
}

func TestAccPrivateLinkService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicIp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.primary").HasValue("true"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicIp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateLinkService_move(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.moveSetup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.2.17"),
			),
		},
		data.ImportStep(),
		{
			Config: r.moveAdd(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.2.17"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.2.18"),
				check.That(data.ResourceName).Key("nat_ip_configuration.2.private_ip_address").HasValue("10.5.2.19"),
				check.That(data.ResourceName).Key("nat_ip_configuration.3.private_ip_address").HasValue("10.5.2.20"),
			),
		},
		data.ImportStep(),
		{
			Config: r.moveChangeOne(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.2.17"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.2.18"),
				check.That(data.ResourceName).Key("nat_ip_configuration.2.private_ip_address").HasValue("10.5.2.19"),
				check.That(data.ResourceName).Key("nat_ip_configuration.3.private_ip_address").HasValue("10.5.2.21"),
			),
		},
		data.ImportStep(),
		{
			Config: r.moveChangeTwo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.2.17"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.2.20"),
				check.That(data.ResourceName).Key("nat_ip_configuration.2.private_ip_address").HasValue("10.5.2.19"),
				check.That(data.ResourceName).Key("nat_ip_configuration.3.private_ip_address").HasValue("10.5.2.21"),
			),
		},
		data.ImportStep(),
		{
			Config: r.moveChangeThree(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.2.17"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.2.20"),
				check.That(data.ResourceName).Key("nat_ip_configuration.2.private_ip_address").HasValue("10.5.2.19"),
				check.That(data.ResourceName).Key("nat_ip_configuration.3.private_ip_address").HasValue("10.5.2.18"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateLinkService_enableProxyProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enable
			Config: r.enableProxyProtocol(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disable
			Config: r.enableProxyProtocol(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enable
			Config: r.enableProxyProtocol(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateLinkService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_link_service", "test")
	r := PrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_approval_subscription_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("visibility_subscription_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("nat_ip_configuration.#").HasValue("2"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address").HasValue("10.5.1.40"),
				check.That(data.ResourceName).Key("nat_ip_configuration.0.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address").HasValue("10.5.1.41"),
				check.That(data.ResourceName).Key("nat_ip_configuration.1.private_ip_address_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("load_balancer_frontend_ip_configuration_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func (t PrivateLinkServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateLinkServices"]

	resp, err := clients.Network.PrivateLinkServiceClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Private Link Service (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r PrivateLinkServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-basic-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) basicIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-update-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                = "acctestPLS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.30"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_link_service" "import" {
  name                = azurerm_private_link_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, r.basic(data), data.RandomInteger)
}

func (r PrivateLinkServiceResource) enableProxyProtocol(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-basic-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.4.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                  = "acctestPLS-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  enable_proxy_protocol = %t

  nat_ip_configuration {
    name      = "primaryIpConfiguration-%d"
    subnet_id = azurerm_subnet.test.id
    primary   = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, enabled, data.RandomInteger)
}

func (r PrivateLinkServiceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-update-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.3.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.30"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.22"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.23"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.3.24"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) moveSetup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-move-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) moveAdd(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-move-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) moveChangeOne(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-move-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.21"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) moveChangeTwo(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-move-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.21"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) moveChangeThree(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-move-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.2.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.17"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.20"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "thirdaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.19"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  nat_ip_configuration {
    name                       = "fourtharyIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.2.18"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateLinkServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-complete-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.5.1.0/24"

  enforce_private_link_service_network_policies = true
}

resource "azurerm_private_link_service" "test" {
  name                           = "acctestPLS-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.1.40"
    private_ip_address_version = "IPv4"
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration-%d"
    subnet_id                  = azurerm_subnet.test.id
    private_ip_address         = "10.5.1.41"
    private_ip_address_version = "IPv4"
    primary                    = false
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.test.frontend_ip_configuration.0.id
  ]

  tags = {
    env = "test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PrivateLinkServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatelinkservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.test.name
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
