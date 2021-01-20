package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ContainerRegistryResource struct {
}

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
		_, errors := containers.ValidateContainerRegistryName(tc.Value, "azurerm_container_registry")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccContainerRegistry_basic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicManaged(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicManaged(data, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_basic_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicManaged(data, "Premium"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_basic_basic2Premium2basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		{
			Config: r.basicManaged(data, "Premium"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Premium"),
			),
		},
		{
			Config: r.basic_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
	})
}

func TestAccContainerRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccContainerRegistry_geoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	skuPremium := "Premium"
	skuBasic := "Basic"

	data.ResourceTest(t, r, []resource.TestStep{
		// first config creates an ACR with locations
		{
			// TODO: fix this to use dynamic locations
			Config: r.geoReplication(data, skuPremium, `eastus", "westus`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue(skuPremium),
				check.That(data.ResourceName).Key("georeplication_locations.#").HasValue("2"),
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
			),
		},
		// second config udpates the ACR with updated locations
		{
			Config: r.geoReplication(data, skuPremium, `centralus", "eastus`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue(skuPremium),
				check.That(data.ResourceName).Key("georeplication_locations.#").HasValue("2"),
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"centralus"`}),
			),
		},
		// third config udpates the ACR with no location
		{
			Config: r.geoReplicationUpdateWithNoLocation(data, skuPremium),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue(skuPremium),
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckContainerRegistryGeoreplications(data.ResourceName, skuPremium, nil),
			),
		},
		// fourth config updates an ACR with replicas
		{
			Config: r.geoReplication(data, skuPremium, `eastus", "westus`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue(skuPremium),
				check.That(data.ResourceName).Key("georeplication_locations.#").HasValue("2"),
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
			),
		},
		// fifth config updates the SKU to basic and no replicas (should remove the existing replicas if any)
		{
			Config: r.geoReplicationUpdateWithNoLocation_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue(skuBasic),
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckContainerRegistryGeoreplications(data.ResourceName, skuBasic, nil),
			),
		},
	})
}

func TestAccContainerRegistry_networkAccessProfileIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkAccessProfile_ip(data, "Premium"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistry_networkAccessProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")
	r := ContainerRegistryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicManaged(data, "Premium"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.networkAccessProfile_ip(data, "Premium"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.ip_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfile_vnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAccessProfile_both(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkAccessProfile_vnet(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.policies(data, 10),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("0"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("10"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("true"),
			),
		},
		{
			Config: r.policies(data, 20),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.0.default_action").HasValue("Allow"),
				check.That(data.ResourceName).Key("network_rule_set.0.virtual_network.#").HasValue("0"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("20"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("true"),
			),
		},
		{
			Config: r.policies_downgradeUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_rule_set.#").HasValue("0"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("trust_policy.0.enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (t ContainerRegistryResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["registries"]

	resp, err := clients.Containers.RegistriesClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Container Registry (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckContainerRegistryGeoreplications(resourceName string, sku string, expectedLocations []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.ReplicationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		resp, err := conn.List(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on containerRegistryClient: %+v", err)
		}

		georeplicationValues := resp.Values()
		expectedLocationsCount := len(expectedLocations) + 1 // the main location is returned by the API as a geolocation for replication.

		// if Sku is not premium, listing the geo-replications locations returns an empty list
		if !strings.EqualFold(sku, string(containerregistry.Premium)) {
			expectedLocationsCount = 0
		}

		actualLocationsCount := len(georeplicationValues)

		if expectedLocationsCount != actualLocationsCount {
			return fmt.Errorf("Bad: Container Registry %q (resource group: %q) expected locations count is %d, actual location count is %d", name, resourceGroup, expectedLocationsCount, actualLocationsCount)
		}

		return nil
	}
}

func (ContainerRegistryResource) basic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automaticcally populate network_rule_set.default_action to allow
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
  name     = "acctestRG-aks-%d"
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
  name     = "acctestRG-aks-%d"
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
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = true
  sku                 = "Basic"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ContainerRegistryResource) geoReplication(data acceptance.TestData, sku string, georeplicationLocations string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                     = "testacccr%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku                      = "%s"
  georeplication_locations = ["%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku, georeplicationLocations)
}

func (ContainerRegistryResource) geoReplicationUpdateWithNoLocation(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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

func (ContainerRegistryResource) geoReplicationUpdateWithNoLocation_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automaticcally populate network_rule_set.default_action to allow
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
  address_prefix       = "10.0.1.0/24"

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
  address_prefix       = "10.0.1.0/24"

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

  retention_policy {
    days    = %d
    enabled = true
  }

  trust_policy {
    enabled = true
  }

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
