// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiManagementDataSource struct{}

func TestAccDataSourceApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("tenant_access.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceApiManagement_tenantAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tenantAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tenant_access.0.enabled").Exists(),
			),
		},
	})
}

func TestAccDataSourceApiManagement_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
	})
}

func TestAccDataSourceApiManagement_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
	})
}

func TestAccDataSourceApiManagement_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management", "test")
	r := ApiManagementDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("publisher_email").HasValue("pub1@email.com"),
				check.That(data.ResourceName).Key("publisher_name").HasValue("pub1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("private_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.public_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("additional_location.0.private_ip_addresses.#").Exists(),
			),
		},
	})
}

func (ApiManagementDataSource) tenantAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%[1]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%[1]d@example.com"
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}

resource "azurerm_api_management_subscription" "test" {
  subscription_id     = "This-Is-A-Valid-Subscription-ID"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  user_id             = azurerm_api_management_user.test.id
  product_id          = azurerm_api_management_product.test.id
  display_name        = "Butter Parser API Enterprise Edition"
  state               = "active"
  allow_tracing       = false
  primary_key         = data.azurerm_api_management.test.tenant_access[0].primary_key
  secondary_key       = data.azurerm_api_management.test.tenant_access[0].secondary_key
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ApiManagementDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementDataSource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementDataSource) identityUserAssigned(data acceptance.TestData) string {
	template := ApiManagementResource{}.identityUserAssigned(data)
	return fmt.Sprintf(`
%s

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, template)
}

func (ApiManagementDataSource) virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestVNET1-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestSNET1-%d"
  resource_group_name  = azurerm_resource_group.test1.name
  virtual_network_name = azurerm_virtual_network.test1.name
  address_prefixes     = ["10.0.1.0/24"]
}

# define internal security_rule according to https://learn.microsoft.com/en-us/azure/api-management/api-management-using-with-internal-vnet#configure-nsg-rules
resource "azurerm_network_security_group" "test1" {
  name                = "acctest-nsg1-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name

  security_rule {
    name                       = "Allow-APIM-Management-Inbound-3443"
    priority                   = 210
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3443"
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
    description                = "Management endpoint for Azure portal and PowerShell"
  }

  security_rule {
    name                       = "Allow-AzureLoadBalancer-Inbound-6390"
    priority                   = 220
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "6390"
    source_address_prefix      = "AzureLoadBalancer"
    destination_address_prefix = "VirtualNetwork"
    description                = "Azure Infrastructure Load Balancer"
  }

  security_rule {
    name                       = "Allow-Internet-Outbound-80"
    priority                   = 240
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Internet"
    description                = "Certificate validation/management"
  }

  security_rule {
    name                       = "Allow-Storage-Outbound-443"
    priority                   = 250
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Storage"
    description                = "Azure Storage dependency (core)"
  }

  security_rule {
    name                       = "Allow-SQL-Outbound-1433"
    priority                   = 260
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "1433"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Sql"
    description                = "Azure SQL dependency (core)"
  }

  security_rule {
    name                       = "Allow-KeyVault-Outbound-443"
    priority                   = 270
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "AzureKeyVault"
    description                = "Azure Key Vault dependency (core)"
  }

  security_rule {
    name                       = "Allow-AzureMonitor-Outbound-1886-443"
    priority                   = 280
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = ["1886", "443"]
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "AzureMonitor"
    description                = "Azure Monitor dependency for diagnostics/metrics/health/App Insights"
  }
}

resource "azurerm_subnet_network_security_group_association" "test1" {
  subnet_id                 = azurerm_subnet.test1.id
  network_security_group_id = azurerm_network_security_group.test1.id
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestVNET2-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestSNET2-%d"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.1.1.0/24"]
}

# define external security_rule according to https://learn.microsoft.com/en-us/azure/api-management/api-management-using-with-internal-vnet#configure-nsg-rules
resource "azurerm_network_security_group" "test2" {
  name                = "acctest-nsg2-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  security_rule {
    name                       = "Allow-APIM-Client-Inbound-80-443"
    priority                   = 200
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = ["80", "443"]
    source_address_prefix      = "Internet"
    destination_address_prefix = "VirtualNetwork"
    description                = "Client communication to API Management (external mode only)"
  }

  security_rule {
    name                       = "Allow-APIM-Management-Inbound-3443"
    priority                   = 210
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3443"
    source_address_prefix      = "ApiManagement"
    destination_address_prefix = "VirtualNetwork"
    description                = "Management endpoint for Azure portal and PowerShell"
  }

  security_rule {
    name                       = "Allow-AzureLoadBalancer-Inbound-6390"
    priority                   = 220
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "6390"
    source_address_prefix      = "AzureLoadBalancer"
    destination_address_prefix = "VirtualNetwork"
    description                = "Azure Infrastructure Load Balancer"
  }

  security_rule {
    name                       = "Allow-AzureTrafficManager-Inbound-443"
    priority                   = 230
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "AzureTrafficManager"
    destination_address_prefix = "VirtualNetwork"
    description                = "Azure Traffic Manager routing for multi-region deployment (external mode only)"
  }

  security_rule {
    name                       = "Allow-Internet-Outbound-80"
    priority                   = 240
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Internet"
    description                = "Validation and management of Microsoft-managed and customer-managed certificates"
  }

  security_rule {
    name                       = "Allow-Storage-Outbound-443"
    priority                   = 250
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Storage"
    description                = "Dependency on Azure Storage for core service functionality"
  }

  security_rule {
    name                       = "Allow-SQL-Outbound-1433"
    priority                   = 260
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "1433"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "Sql"
    description                = "Access to Azure SQL endpoints for core service functionality"
  }

  security_rule {
    name                       = "Allow-KeyVault-Outbound-443"
    priority                   = 270
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "AzureKeyVault"
    description                = "Access to Azure Key Vault for core service functionality"
  }

  security_rule {
    name                       = "Allow-AzureMonitor-Outbound-1886-443"
    priority                   = 280
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = ["1886", "443"]
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "AzureMonitor"
    description                = "Publish diagnostics logs and metrics, Resource Health, and Application Insights"
  }
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"

  additional_location {
    location = azurerm_resource_group.test2.location
    virtual_network_configuration {
      subnet_id = azurerm_subnet.test2.id
    }
  }

  virtual_network_type = "Internal"
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test1.id
  }

  depends_on = [
    azurerm_subnet_network_security_group_association.test1,
    azurerm_subnet_network_security_group_association.test2
  ]
}

data "azurerm_api_management" "test" {
  name                = azurerm_api_management.test.name
  resource_group_name = azurerm_api_management.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
