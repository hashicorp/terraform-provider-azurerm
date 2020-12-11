package containers_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerRegistryName_validation(t *testing.T) {
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
		_, errors := ValidateAzureRMContainerRegistryName(tc.Value, "azurerm_container_registry")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccAzureRMContainerRegistry_basic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMContainerRegistry_requiresImport(data, "Basic"),
				ExpectError: acceptance.RequiresImportError("azurerm_container_registry"),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_basic2Premium2basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Premium"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_completeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_geoReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	skuPremium := "Premium"
	skuBasic := "Basic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			// first config creates an ACR with locations
			{
				// TODO: fix this to use dynamic locations
				Config: testAccAzureRMContainerRegistry_geoReplication(data, skuPremium, `eastus", "westus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(data.ResourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					testCheckAzureRMContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// second config udpates the ACR with updated locations
			{
				Config: testAccAzureRMContainerRegistry_geoReplication(data, skuPremium, `centralus", "eastus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(data.ResourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					testCheckAzureRMContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"centralus"`}),
				),
			},
			// third config udpates the ACR with no location
			{
				Config: testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(data, skuPremium),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", skuPremium),
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					testCheckAzureRMContainerRegistryGeoreplications(data.ResourceName, skuPremium, nil),
				),
			},
			// fourth config updates an ACR with replicas
			{
				Config: testAccAzureRMContainerRegistry_geoReplication(data, skuPremium, `eastus", "westus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", skuPremium),
					resource.TestCheckResourceAttr(data.ResourceName, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					testCheckAzureRMContainerRegistryGeoreplications(data.ResourceName, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// fifth config updates the SKU to basic and no replicas (should remove the existing replicas if any)
			{
				Config: testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", skuBasic),
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					testCheckAzureRMContainerRegistryGeoreplications(data.ResourceName, skuBasic, nil),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_networkAccessProfileIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_ip(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.ip_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_networkAccessProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_ip(data, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.ip_rule.#", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.virtual_network.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_both(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.ip_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_networkAccessProfileVnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_vnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.virtual_network.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistry_policies(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_policies(data, 10),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.virtual_network.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "trust_policy.0.enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_policies(data, 20),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.virtual_network.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "trust_policy.0.enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_policies_downgradeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "trust_policy.0.enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMContainerRegistryDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.RegistriesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_registry" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testCheckAzureRMContainerRegistryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.RegistriesClient
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

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on containerRegistryClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Container Registry %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMContainerRegistryGeoreplications(resourceName string, sku string, expectedLocations []string) resource.TestCheckFunc {
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

func testAccAzureRMContainerRegistry_basic_basic(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_basicManaged(data acceptance.TestData, sku string) string {
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

func testAccAzureRMContainerRegistry_requiresImport(data acceptance.TestData, sku string) string {
	template := testAccAzureRMContainerRegistry_basicManaged(data, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "import" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
  location            = azurerm_container_registry.test.location
  sku                 = azurerm_container_registry.test.sku
}
`, template)
}

func testAccAzureRMContainerRegistry_complete(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_completeUpdated(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_geoReplication(data acceptance.TestData, sku string, georeplicationLocations string) string {
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

func testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(data acceptance.TestData, sku string) string {
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

func testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation_basic(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_networkAccessProfile_ip(data acceptance.TestData, sku string) string {
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

func testAccAzureRMContainerRegistry_networkAccessProfile_vnet(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_networkAccessProfile_both(data acceptance.TestData) string {
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

func testAccAzureRMContainerRegistry_policies(data acceptance.TestData, days int) string {
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

func testAccAzureRMContainerRegistry_policies_downgradeUpdate(data acceptance.TestData) string {
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
