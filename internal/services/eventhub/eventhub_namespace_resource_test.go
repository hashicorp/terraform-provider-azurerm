// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EventHubNamespaceResource struct{}

func TestAccEventHubNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

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

func TestAccEventHubNamespace_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_basicUpdateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventhub_namespace"),
		},
	})
}

func TestAccEventHubNamespace_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_standardWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_standardUpdateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_iprule_trusted_services(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkrule_iprule_trusted_services(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_iprule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkrule_iprule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_publicNetworkAccessDiff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.networkrule_publicNetworkAccessDiff(data),
			ExpectError: regexp.MustCompile("the value of public network access of namespace should be the same as of the network rulesets"),
		},
	})
}

func TestAccEventHubNamespace_networkrule_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkrule_vnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkruleVnetIpRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkruleVnetIpRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rulesets.0.virtual_network_rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("network_rulesets.0.ip_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_readDefaultKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_primary_key", regexp.MustCompile(".+")),
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_secondary_key", regexp.MustCompile(".+")),
			),
		},
	})
}

func TestAccEventHubNamespace_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// `default_primary_connection_string_alias` and `default_secondary_connection_string_alias` are still `nil` in `azurerm_eventhub_namespace` after created `azurerm_eventhub_namespace` successfully since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
			// So these two properties should be checked in the second run.
			Config: r.withAliasConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withAliasConnectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_primary_connection_string_alias", regexp.MustCompile("Endpoint=.+")),
				acceptance.TestMatchResourceAttr(data.ResourceName, "default_secondary_connection_string_alias", regexp.MustCompile("Endpoint=.+")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_maximumThroughputUnits(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maximumThroughputUnits(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_dedicatedClusterID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedClusterID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_BasicWithTagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicWithTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("3"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.sauce").HasValue("Hot"),
				check.That(data.ResourceName).Key("tags.terraform").HasValue("true"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacity(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("20"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithLocalAuthProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.localAuthProperty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithCapacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacity(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("20"),
			),
		},
		{
			Config: r.capacity(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithSkuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		data.ImportStep(),
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_SkuDowngradeFromAutoInflateWithMaxThroughput(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maximumThroughputUnits(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
	})
}

func TestAccEventHubNamespace_maximumThroughputUnitsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maximumThroughputUnits(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
				check.That(data.ResourceName).Key("maximum_throughput_units").HasValue("25"),
			),
		},
		{
			Config: r.maximumThroughputUnitsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("1"),
				check.That(data.ResourceName).Key("maximum_throughput_units").HasValue("1"),
			),
		},
	})
}

func TestAccEventHubNamespace_publicNetworkAccessUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

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

func TestAccEventHubNamespace_minimumTLSUpdate(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("The `minimum_tls_version` has only one possible value `1.2`, we can not update it.")
	}
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

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

func TestAccEventHubNamespace_autoInfalteDisabledWithAutoInflateUnits(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoInfalteDisabledWithAutoInflateUnits(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (EventHubNamespaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := namespaces.ParseNamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.NamespacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (EventHubNamespaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) withAliasConnectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-eh-%[1]d-2"
  location = "%[3]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace" "test2" {
  name                = "acctesteventhubnamespace2-%[1]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  sku = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  namespace_name       = azurerm_eventhub_namespace.test.name
  partner_namespace_id = azurerm_eventhub_namespace.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (EventHubNamespaceResource) requiresImport(data acceptance.TestData) string {
	template := EventHubNamespaceResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace" "import" {
  name                = azurerm_eventhub_namespace.test.name
  location            = azurerm_eventhub_namespace.test.location
  resource_group_name = azurerm_eventhub_namespace.test.resource_group_name
  sku                 = azurerm_eventhub_namespace.test.sku
}
`, template)
}

func (EventHubNamespaceResource) standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"
  network_rulesets {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) standardWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) networkrule_iprule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"

  network_rulesets {
    default_action = "Deny"
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) networkrule_publicNetworkAccessDiff(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                          = "acctesteventhubnamespace-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  sku                           = "Standard"
  capacity                      = "2"
  public_network_access_enabled = true

  network_rulesets {
    default_action                = "Deny"
    public_network_access_enabled = false
    ip_rule {
      ip_mask = "10.0.0.0/16"
      action  = "Allow"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) networkrule_iprule_trusted_services(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"

  network_rulesets {
    default_action                 = "Deny"
    trusted_service_access_enabled = true
    ip_rule {
      ip_mask = "10.0.0.0/16"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) networkrule_vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"

  network_rulesets {
    default_action = "Deny"
    virtual_network_rule {
      subnet_id = azurerm_subnet.test.id

      ignore_missing_virtual_network_service_endpoint = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventHubNamespaceResource) networkruleVnetIpRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn1-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub1-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.EventHub"]
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctvn2-%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctsub2-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.1.1.0/24"]
  service_endpoints    = ["Microsoft.EventHub"]
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"

  network_rulesets {
    default_action = "Deny"

    virtual_network_rule {
      subnet_id = azurerm_subnet.test.id
    }

    virtual_network_rule {
      subnet_id = azurerm_subnet.test2.id
    }

    ip_rule {
      ip_mask = "10.0.1.0/24"
    }

    ip_rule {
      ip_mask = "10.1.1.0/24"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventHubNamespaceResource) maximumThroughputUnits(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "Standard"
  capacity                 = "2"
  auto_inflate_enabled     = true
  maximum_throughput_units = 25
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) zoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) dedicatedClusterID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctesteventhubcluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated_1"
}

resource "azurerm_eventhub_namespace" "test" {
  name                 = "acctesteventhubnamespace-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku                  = "Standard"
  capacity             = "2"
  dedicated_cluster_id = azurerm_eventhub_cluster.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (EventHubNamespaceResource) basicWithTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  tags = {
    environment = "Production"
    sauce       = "Hot"
    terraform   = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) capacity(data acceptance.TestData, capacity int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  capacity            = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacity)
}

func (EventHubNamespaceResource) localAuthProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                         = "acctesteventhubnamespace-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "Basic"
  local_authentication_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) publicNetworkAccessUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                          = "acctesteventhubnamespace-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  sku                           = "Basic"
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) minimumTLSUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
  minimum_tls_version = "1.1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) maximumThroughputUnitsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "Standard"
  capacity                 = 1
  auto_inflate_enabled     = true
  maximum_throughput_units = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) autoInfalteDisabledWithAutoInflateUnits(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eh-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "Standard"
  capacity                 = 1
  auto_inflate_enabled     = false
  maximum_throughput_units = 0
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
