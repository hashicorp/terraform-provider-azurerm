// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryResource struct{}

func TestAccContainerRegistryName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "four",
			ErrCount: 1,
		},
		{
			Value:    "5five",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},

		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd33241202",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd3324120",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validate.ContainerRegistryName(tc.Value, "azurerm_container_registry")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccContainerRegistry_basic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

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

func TestAccContainerRegistry_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManaged(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, "Basic"),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry"),
		},
	})
}

func TestAccContainerRegistry_basic_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManaged(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_basic_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManaged(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_basic_basic2Premium2basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		{
			Config: r.basicManaged(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Premium"),
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

func TestAccContainerRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

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

func TestAccContainerRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccContainerRegistry_geoReplicationLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	secondaryLocation := location.Normalize(data.Locations.Secondary)
	ternaryLocation := location.Normalize(data.Locations.Ternary)

	data.ResourceTest(t, r, []acceptance.TestStep{
		// creates an ACR with locations
		{
			Config: r.geoReplicationLocation(data, secondaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the ACR with updated locations
		{
			Config: r.geoReplicationLocation(data, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the ACR with updated locations
		{
			Config: r.geoReplicationMultipleLocations(data, secondaryLocation, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoReplicationMultipleLocationsUpdate(data, secondaryLocation, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the SKU to basic.
		{
			Config: r.geoReplicationUpdateWithNoLocation_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_geoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	secondaryLocation := location.Normalize(data.Locations.Secondary)
	ternaryLocation := location.Normalize(data.Locations.Ternary)

	data.ResourceTest(t, r, []acceptance.TestStep{
		// creates an ACR with locations
		{
			Config: r.geoReplication(data, secondaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the ACR with updated locations
		{
			Config: r.geoReplication(data, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the ACR with updated locations
		{
			Config: r.geoReplicationMultipleLocations(data, secondaryLocation, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoReplicationMultipleLocationsUpdate(data, secondaryLocation, ternaryLocation),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// updates the SKU to basic.
		{
			Config: r.geoReplicationUpdateWithNoReplication_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkAccessProfileIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAccessProfile_ip(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkAccessProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicManaged(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.networkAccessProfile_ip(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfile_vnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfile_both(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkAccessProfileVnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkAccessProfile_vnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_policies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policies(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("0"),
				check.That(data.ResourceName).Key("quarantine_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("10"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("export_policy_enabled").HasValue("false"),
			),
		},
		{
			Config: r.policies(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("0"),
				check.That(data.ResourceName).Key("quarantine_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("20"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("export_policy_enabled").HasValue("false"),
			),
		},
		{
			Config: r.policies_downgradeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("quarantine_policy_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("export_policy_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// creates an ACR with encryption
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Premium"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// creates an ACR with encryption
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Premium"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_zoneRedundancy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundancy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_geoReplicationZoneRedundancy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoReplicationZoneRedundancy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_geoReplicationRegionEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.regionEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_anonymousPull(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.anonymousPullStandard(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.anonymousPullStandard(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.anonymousPullStandard(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_dataEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataEndpointPremium(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataEndpointPremium(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dataEndpointPremium(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkRuleBypassOption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkRuleBypassOptionsPremium(data, "None"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkRuleBypassOptionsPremium(data, "AzureServices"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkRuleBypassOptionsPremium(data, "None"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ContainerRegistryResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := registries.ParseRegistryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.ContainerRegistryClient_v2021_08_01_preview.Registries.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ContainerRegistryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automatically populate network_rule_set.default_action to allow
  network_rule_set = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) basicManaged(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func (r ContainerRegistryResource) requiresImport(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "import" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
  location            = azurerm_container_registry.test.location
  sku                 = azurerm_container_registry.test.sku
}
`, r.basicManaged(data, sku))
}

func (ContainerRegistryResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Basic"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) completeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Premium"

  tags = {
    environment = "production"
  }
  public_network_access_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) geoReplicationLocation(data acceptance.TestData, replicationRegion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, replicationRegion)
}

func (ContainerRegistryResource) geoReplication(data acceptance.TestData, replication string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location = "%s"
    tags = {
      Environment = "Production"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, replication)
}

func (ContainerRegistryResource) geoReplicationMultipleLocations(data acceptance.TestData, primaryLocation string, secondaryLocation string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location = "%s"
  }
  georeplications {
    location = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, primaryLocation, secondaryLocation)
}

func (ContainerRegistryResource) geoReplicationMultipleLocationsUpdate(data acceptance.TestData, primaryLocation string, secondaryLocation string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location                = "%s"
    zone_redundancy_enabled = true
  }
  georeplications {
    location                  = "%s"
    regional_endpoint_enabled = true
    tags = {
      foo = "bar"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, primaryLocation, secondaryLocation)
}

func (ContainerRegistryResource) geoReplicationUpdateWithNoLocation_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premium SKU will automatically populate network_rule_set.default_action to allow
  network_rule_set = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) geoReplicationUpdateWithNoReplication_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automatically populate network_rule_set.default_action to allow
  network_rule_set = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) networkAccessProfile_ip(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testAccCr%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "%[3]s"
  admin_enabled       = false

  network_rule_set {
    default_action = "Allow"

    ip_rule {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }

    ip_rule {
      action   = "Allow"
      ip_range = "1.1.1.1/32"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (ContainerRegistryResource) networkAccessProfile_vnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "virtualNetwork1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  service_endpoints = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_container_registry" "test" {
  name                = "testAccCr%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false

  network_rule_set {
    default_action = "Deny"

    ip_rule {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }

    virtual_network {
      action    = "Allow"
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryResource) networkAccessProfile_both(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "virtualNetwork1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  service_endpoints = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_container_registry" "test" {
  name                = "testAccCr%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false

  network_rule_set {
    default_action = "Deny"

    ip_rule {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }

    virtual_network {
      action    = "Allow"
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryResource) policies(data acceptance.TestData, days int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "acctestACR%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Premium"

  quarantine_policy_enabled = true

  retention_policy {
    days    = %d
    enabled = true
  }

  trust_policy {
    enabled = true
  }

  export_policy_enabled         = false
  public_network_access_enabled = false

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, days)
}

func (ContainerRegistryResource) policies_downgradeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "acctestACR%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Basic"
  network_rule_set    = []

  retention_policy {}
  trust_policy {}

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "testaccuai%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ContainerRegistryResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  identity {
    type = "SystemAssigned"
  }
}


`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) zoneRedundancy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}
resource "azurerm_container_registry" "test" {
  name                    = "testacccr%d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  sku                     = "Premium"
  zone_redundancy_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) geoReplicationZoneRedundancy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}
resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location                = "%s"
    zone_redundancy_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (ContainerRegistryResource) regionEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}
resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  georeplications {
    location                  = "%s"
    regional_endpoint_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (ContainerRegistryResource) anonymousPullStandard(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                   = "testacccr%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  sku                    = "Standard"
  anonymous_pull_enabled = "%t"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (ContainerRegistryResource) dataEndpointPremium(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                  = "testacccr%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  sku                   = "Premium"
  data_endpoint_enabled = "%t"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (ContainerRegistryResource) networkRuleBypassOptionsPremium(data acceptance.TestData, opt string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                       = "testacccr%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku                        = "Premium"
  network_rule_bypass_option = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, opt)
}
