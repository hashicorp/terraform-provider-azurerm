package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
		_, errors := validateAzureRMContainerRegistryName(tc.Value, "azurerm_container_registry")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Container Registry Name to trigger a validation error: %v", errors)
		}
	}
}

func TestAccAzureRMContainerRegistry_basic_basic(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(ri, l, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				Config:      testAccAzureRMContainerRegistry_requiresImport(ri, l, "Basic"),
				ExpectError: testRequiresImportError("azurerm_container_registry"),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_standard(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_premium(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_basic_basic2Premium2basic(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
					resource.TestCheckResourceAttr(rn, "sku", "Basic"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(ri, testLocation(), "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
					resource.TestCheckResourceAttr(rn, "sku", "Premium"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
					resource.TestCheckResourceAttr(rn, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_complete(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_update(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_complete(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_completeUpdated(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_geoReplication(t *testing.T) {
	dsn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()

	skuPremium := "Premium"
	skuBasic := "Basic"

	containerRegistryName := fmt.Sprintf("testacccr%d", ri)
	resourceGroupName := fmt.Sprintf("testAccRg-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			// first config creates an ACR with locations
			{
				Config: testAccAzureRMContainerRegistry_geoReplication(ri, testLocation(), skuPremium, `eastus", "westus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dsn, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dsn, "sku", skuPremium),
					resource.TestCheckResourceAttr(dsn, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dsn),
					testCheckAzureRMContainerRegistryGeoreplications(dsn, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// second config udpates the ACR with updated locations
			{
				Config: testAccAzureRMContainerRegistry_geoReplication(ri, testLocation(), skuPremium, `centralus", "eastus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dsn, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dsn, "sku", skuPremium),
					resource.TestCheckResourceAttr(dsn, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dsn),
					testCheckAzureRMContainerRegistryGeoreplications(dsn, skuPremium, []string{`"eastus"`, `"centralus"`}),
				),
			},
			// third config udpates the ACR with no location
			{
				Config: testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(ri, testLocation(), skuPremium),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dsn, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dsn, "sku", skuPremium),
					testCheckAzureRMContainerRegistryExists(dsn),
					testCheckAzureRMContainerRegistryGeoreplications(dsn, skuPremium, nil),
				),
			},
			// fourth config updates an ACR with replicas
			{
				Config: testAccAzureRMContainerRegistry_geoReplication(ri, testLocation(), skuPremium, `eastus", "westus`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dsn, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dsn, "sku", skuPremium),
					resource.TestCheckResourceAttr(dsn, "georeplication_locations.#", "2"),
					testCheckAzureRMContainerRegistryExists(dsn),
					testCheckAzureRMContainerRegistryGeoreplications(dsn, skuPremium, []string{`"eastus"`, `"westus"`}),
				),
			},
			// fifth config updates the SKU to basic and no replicas (should remove the existing replicas if any)
			{
				Config: testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "name", containerRegistryName),
					resource.TestCheckResourceAttr(dsn, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dsn, "sku", skuBasic),
					testCheckAzureRMContainerRegistryExists(dsn),
					testCheckAzureRMContainerRegistryGeoreplications(dsn, skuBasic, nil),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_networkAccessProfile_ip(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_ip(ri, l, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
					resource.TestCheckResourceAttr(rn, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(rn, "network_rule_set.0.ip_rule.#", "1"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistry_networkAccessProfile_update(t *testing.T) {
	rn := "azurerm_container_registry.test"
	ri := tf.AccRandTimeInt()
	l := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistry_basicManaged(ri, l, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_networkAccessProfile_ip(ri, l, "Premium"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
					resource.TestCheckResourceAttr(rn, "network_rule_set.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(rn, "network_rule_set.0.ip_rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistry_basic_basic(ri, l),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryExists(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMContainerRegistryDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).containers.RegistriesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		conn := testAccProvider.Meta().(*ArmClient).containers.RegistriesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		conn := testAccProvider.Meta().(*ArmClient).containers.ReplicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMContainerRegistry_basic_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automaticcally populate network_rule_set.default_action to allow
  network_rule_set = []
}
`, rInt, location, rInt)
}

func testAccAzureRMContainerRegistry_basicManaged(rInt int, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "%s"
}
`, rInt, location, rInt, sku)
}

func testAccAzureRMContainerRegistry_requiresImport(rInt int, location string, sku string) string {
	template := testAccAzureRMContainerRegistry_basicManaged(rInt, location, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "import" {
  name                = "${azurerm_container_registry.test.name}"
  resource_group_name = "${azurerm_container_registry.test.resource_group_name}"
  location            = "${azurerm_container_registry.test.location}"
  sku                 = "${azurerm_container_registry.test.sku}"

}
`, template)
}

func testAccAzureRMContainerRegistry_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  admin_enabled       = false
  sku                 = "Basic"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMContainerRegistry_completeUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  admin_enabled       = true
  sku                 = "Basic"

  tags = {
    environment = "production"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMContainerRegistry_geoReplication(rInt int, location string, sku string, georeplicationLocations string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                     = "testacccr%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  sku                      = "%s"
  georeplication_locations = ["%s"]
}
`, rInt, location, rInt, sku, georeplicationLocations)
}

func testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation(rInt int, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "%s"
}
`, rInt, location, rInt, sku)
}

func testAccAzureRMContainerRegistry_geoReplicationUpdateWithNoLocation_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "Basic"

  # make sure network_rule_set is empty for basic SKU
  # premiuim SKU will automaticcally populate network_rule_set.default_action to allow
  network_rule_set = []
}
`, rInt, location, rInt)
}

func testAccAzureRMContainerRegistry_networkAccessProfile_ip(rInt int, location string, sku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testAccRg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_container_registry" "test" {
  name                = "testAccCr%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "%[3]s"
  admin_enabled       = false

  network_rule_set {
    default_action = "Allow"

	ip_rule {
      action   = "Allow"
      ip_range = "8.8.8.8/32"
    }
  }
}
`, rInt, location, sku)
}
