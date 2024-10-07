// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusNamespaceResource struct{}

func TestAccAzureRMServiceBusNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.standardSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMServiceBusNamespace_readDefaultKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestMatchResourceAttr(
					data.ResourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
				acceptance.TestMatchResourceAttr(
					data.ResourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
				acceptance.TestMatchResourceAttr(
					data.ResourceName, "default_primary_key", regexp.MustCompile(".+")),
				acceptance.TestMatchResourceAttr(
					data.ResourceName, "default_secondary_key", regexp.MustCompile(".+")),
			),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.premium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_basicCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicCapacity(data),
			ExpectError: regexp.MustCompile("service bus SKU \"Basic\" only supports `capacity` of 0"),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premiumCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.premiumCapacity(data),
			ExpectError: regexp.MustCompile("service bus SKU \"Premium\" only supports `capacity` of 1, 2, 4, 8 or 16"),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_premiumMessagingPartition(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.premiumMessagingPartition(data),
			ExpectError: regexp.MustCompile("service bus SKU \"Premium\" only supports `premium_messaging_partitions` of 1, 2, 4"),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_zoneRedundant(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipped as 'zone_redundant' property is deprecated in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_publicNetworkAccessUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		{
			Config: r.publicNetworkAccessUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_minimumTLSUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("minimum_tls_version").HasValue("1.2"),
			),
		},
		{
			Config: r.minimumTLSUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("minimum_tls_version").HasValue("1.1"),
			),
		},
	})
}

func TestAccAzureRMServiceBusNamespace_endpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestMatchResourceAttr(
					data.ResourceName, "endpoint", regexp.MustCompile(`https://.+`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMServiceBusNamespace_networkRuleSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_namespace", "test")
	r := ServiceBusNamespaceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkRuleSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Deny"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkRuleSetComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.trusted_services_allowed").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkRuleSetEmpty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ServiceBusNamespaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespaces.ParseNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceBus.NamespacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Service Bus NameSpace (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ServiceBusNamespaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  local_auth_enabled  = false
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ServiceBusNamespaceResource) standardSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  local_auth_enabled  = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServiceBusNamespaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace" "import" {
  name                = azurerm_servicebus_namespace.test.name
  location            = azurerm_servicebus_namespace.test.location
  resource_group_name = azurerm_servicebus_namespace.test.resource_group_name
  sku                 = azurerm_servicebus_namespace.test.sku
}
`, r.basic(data))
}

func (ServiceBusNamespaceResource) premium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  capacity                     = 4
  premium_messaging_partitions = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) basicCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  capacity            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) premiumCapacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  capacity                     = 0
  premium_messaging_partitions = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) premiumMessagingPartition(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium"
  capacity            = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

// TODO: Remove in v4.0
func (ServiceBusNamespaceResource) zoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1
  zone_redundant               = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusNamespaceResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ServiceBusNamespaceResource) customerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"
    ]
    secret_permissions = [
      "Get",
    ]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id
    key_permissions = [
      "Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"
    ]
    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%[2]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  premium_messaging_partitions = 1
  capacity                     = 1

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  customer_managed_key {
    key_vault_key_id                  = azurerm_key_vault_key.test.id
    identity_id                       = azurerm_user_assigned_identity.test.id
    infrastructure_encryption_enabled = true
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (ServiceBusNamespaceResource) publicNetworkAccessUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                          = "acctestservicebusnamespace-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  sku                           = "Basic"
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) minimumTLSUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  minimum_tls_version = "1.1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ServiceBusNamespaceResource) networkRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "${azurerm_virtual_network.test.name}-default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  capacity                     = 1
  premium_messaging_partitions = 1

  network_rule_set {
    default_action = "Deny"

    network_rules {
      subnet_id                            = azurerm_subnet.test.id
      ignore_missing_vnet_service_endpoint = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ServiceBusNamespaceResource) networkRuleSetComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "${azurerm_virtual_network.test.name}-default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  capacity                     = 1
  premium_messaging_partitions = 1
  network_rule_set {
    default_action                = "Deny"
    trusted_services_allowed      = true
    public_network_access_enabled = true

    ip_rules = ["1.1.1.1"]

    network_rules {
      subnet_id                            = azurerm_subnet.test.id
      ignore_missing_vnet_service_endpoint = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ServiceBusNamespaceResource) networkRuleSetEmpty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                 = "${azurerm_virtual_network.test.name}-default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.ServiceBus"]
}

resource "azurerm_servicebus_namespace" "test" {
  name                         = "acctestservicebusnamespace-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Premium"
  capacity                     = 1
  premium_messaging_partitions = 1

  network_rule_set {}
}
`, data.RandomInteger, data.Locations.Primary)
}
