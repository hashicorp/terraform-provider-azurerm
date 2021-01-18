package eventhub_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventHubNamespaceResource struct {
}

func TestAccEventHubNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

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

func TestAccEventHubNamespace_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_basicUpdateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_standardWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardWithIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_standardUpdateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardWithIdentity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_iprule_trusted_services(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkrule_iprule_trusted_services(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_iprule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkrule_iprule(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkrule_vnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkrule_vnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_networkruleVnetIpRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkruleVnetIpRule(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestMatchResourceAttr(data.ResourceName, "default_primary_connection_string", regexp.MustCompile("Endpoint=.+")),
				resource.TestMatchResourceAttr(data.ResourceName, "default_secondary_connection_string", regexp.MustCompile("Endpoint=.+")),
				resource.TestMatchResourceAttr(data.ResourceName, "default_primary_key", regexp.MustCompile(".+")),
				resource.TestMatchResourceAttr(data.ResourceName, "default_secondary_key", regexp.MustCompile(".+")),
			),
		},
	})
}

func TestAccEventHubNamespace_withAliasConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// `default_primary_connection_string_alias` and `default_secondary_connection_string_alias` are still `nil` in `azurerm_eventhub_namespace` after created `azurerm_eventhub_namespace` successfully since `azurerm_eventhub_namespace_disaster_recovery_config` hasn't been created.
			// So these two properties should be checked in the second run.
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withAliasConnectionString(data),
			Check: resource.ComposeTestCheckFunc(
				resource.TestMatchResourceAttr(data.ResourceName, "default_primary_connection_string_alias", regexp.MustCompile("Endpoint=.+")),
				resource.TestMatchResourceAttr(data.ResourceName, "default_secondary_connection_string_alias", regexp.MustCompile("Endpoint=.+")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_maximumThroughputUnits(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.maximumThroughputUnits(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_dedicatedClusterID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dedicatedClusterID(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventHubNamespace_NonStandardCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nonStandardCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:             r.nonStandardCasing(data),
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	})
}

func TestAccEventHubNamespace_BasicWithTagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicWithTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithCapacity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.capacity(data, 20),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("20"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithCapacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.capacity(data, 20),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("20"),
			),
		},
		{
			Config: r.capacity(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
			),
		},
	})
}

func TestAccEventHubNamespace_BasicWithSkuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
			),
		},
	})
}

func TestAccEventHubNamespace_maximumThroughputUnitsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.maximumThroughputUnits(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("2"),
				check.That(data.ResourceName).Key("maximum_throughput_units").HasValue("20"),
			),
		},
		{
			Config: r.maximumThroughputUnitsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("capacity").HasValue("1"),
				check.That(data.ResourceName).Key("maximum_throughput_units").HasValue("1"),
			),
		},
	})
}

func TestAccEventHubNamespace_autoInfalteDisabledWithAutoInflateUnits(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace", "test")
	r := EventHubNamespaceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoInfalteDisabledWithAutoInflateUnits(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (EventHubNamespaceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.NamespacesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.EHNamespaceProperties != nil), nil
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
  name     = "acctestRG-ehn-%[1]d"
  location = "%[2]s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-ehn-%[1]d"
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
  address_prefix       = "10.0.2.0/24"
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
  name     = "acctestRG-%[1]d"
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
  address_prefix       = "10.0.1.0/24"
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
  address_prefix       = "10.1.1.0/24"
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

func (EventHubNamespaceResource) nonStandardCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) maximumThroughputUnits(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                     = "acctesteventhubnamespace-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "Standard"
  capacity                 = "2"
  auto_inflate_enabled     = true
  maximum_throughput_units = 20
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) zoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = "2"
  zone_redundant      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventHubNamespaceResource) dedicatedClusterID(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"

  tags = {
    environment = "Production"
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
  name     = "acctestRG-%d"
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

func (EventHubNamespaceResource) maximumThroughputUnitsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
  name     = "acctestRG-%d"
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
