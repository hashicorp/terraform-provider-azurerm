// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryResource struct{}

func TestAccContainerRegistry_basic(t *testing.T) {
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
			Config:      r.requiresImport(data, "Basic"),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry"),
		},
	})
}

func TestAccContainerRegistry_basicManagedStandard(t *testing.T) {
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

func TestAccContainerRegistry_basicManagedPremium(t *testing.T) {
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

func TestAccContainerRegistry_basic2Premium2basic(t *testing.T) {
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
		data.ImportStep(),
		{
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.downgradeSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_geoReplicationLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	locs := []string{location.Normalize(data.Locations.Secondary), location.Normalize(data.Locations.Ternary)}
	// Sorting the secondary and ternary locations to ensure the order as is expected by this resource (see its Read() function)
	slices.Sort(locs)
	secondaryLocation := locs[0]
	ternaryLocation := locs[1]

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
			Config: r.geoReplicationUpdateWithNoLocationBasic(data),
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
			Config: r.networkAccessProfileIp(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfileIpRemoved(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfileNetworkRuleSetRemoved(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkAccessProfileUpdate(t *testing.T) {
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
		{
			Config: r.networkAccessProfileIp(data, "Premium"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_policies(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("Skipping in 4.0 since policy updates are tested in the update test using the new properties")
	}
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policies(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.policies(data, 20),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.policies_downgradeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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

func TestAccContainerRegistry_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionEnabled(data),
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

	resp, err := clients.Containers.ContainerRegistryClient_v2023_06_01_preview.Registries.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ContainerRegistryResource) basic(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
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
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Premium"
  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = false
  quarantine_policy_enabled     = true
  retention_policy {
    enabled = true
    days    = 10
  }
  trust_policy {
    enabled = true
  }
  export_policy_enabled  = false
  anonymous_pull_enabled = true
  data_endpoint_enabled  = true

  network_rule_bypass_option = "None"

  tags = {
    environment = "production"
  }
}
`, data.Locations.Primary, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Premium"
  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = false
  quarantine_policy_enabled     = true
  retention_policy_in_days      = 10
  trust_policy_enabled          = true
  export_policy_enabled         = false
  anonymous_pull_enabled        = true
  data_endpoint_enabled         = true

  network_rule_bypass_option = "None"

  tags = {
    environment = "production"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) completeUpdated(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "testaccuai%[2]d"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Premium"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  public_network_access_enabled = true
  quarantine_policy_enabled     = false
  retention_policy {
    enabled = true
    days    = 15
  }
  trust_policy {
    enabled = false
  }
  export_policy_enabled  = true
  anonymous_pull_enabled = false
  data_endpoint_enabled  = false

  network_rule_bypass_option = "AzureServices"

  tags = {
    environment = "production"
    oompa       = "loompa"
  }
}
`, data.Locations.Primary, data.RandomInteger)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "testaccuai%[2]d"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Premium"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  public_network_access_enabled = true
  quarantine_policy_enabled     = false
  retention_policy_in_days      = 15
  trust_policy_enabled          = false
  export_policy_enabled         = true
  anonymous_pull_enabled        = false
  data_endpoint_enabled         = false

  network_rule_bypass_option = "AzureServices"

  tags = {
    environment = "production"
    oompa       = "loompa"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) downgradeSku(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Basic"

  identity {
    type = "SystemAssigned"
  }

  retention_policy {}
  trust_policy {}

  network_rule_set = []
}
`, data.Locations.Primary, data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Basic"

  identity {
    type = "SystemAssigned"
  }
}
`, data.Locations.Primary, data.RandomInteger)
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

func (ContainerRegistryResource) geoReplicationUpdateWithNoLocationBasic(data acceptance.TestData) string {
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

func (ContainerRegistryResource) networkAccessProfileIp(data acceptance.TestData, sku string) string {
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

func (ContainerRegistryResource) networkAccessProfileIpRemoved(data acceptance.TestData, sku string) string {
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
  }
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (ContainerRegistryResource) networkAccessProfileNetworkRuleSetRemoved(data acceptance.TestData, sku string) string {
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
}
`, data.RandomInteger, data.Locations.Primary, sku)
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

func (ContainerRegistryResource) encryptionEnabled(data acceptance.TestData) string {
	template := ContainerRegistryResource{}.encryptionTemplate(data)

	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  encryption {
    enabled            = true
    identity_client_id = azurerm_user_assigned_identity.test.client_id
    key_vault_key_id   = azurerm_key_vault_key.test.id
  }
}
`, template, data.RandomInteger, data.RandomString)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  encryption {
    identity_client_id = azurerm_user_assigned_identity.test.client_id
    key_vault_key_id   = azurerm_key_vault_key.test.id
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (ContainerRegistryResource) encryptionTemplate(data acceptance.TestData) string {
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
  name     = "acctestRG-acr-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "%[3]s"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
