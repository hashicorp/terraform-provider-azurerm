package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AdvancedThreatProtectionResource struct {
}

func TestAccAdvancedThreatProtection_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advanced_threat_protection", "test")
	r := AdvancedThreatProtectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageAccount(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.storageAccount(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAdvancedThreatProtection_cosmosAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advanced_threat_protection", "test")
	r := AdvancedThreatProtectionResource{}

	// the API errors on deleting the cosmos DB account some of the time so lets skip this test for now
	// TODO: remove once this is fixed: https://github.com/Azure/azure-sdk-for-go/issues/6310
	// run it multiple times in a row as it only fails 50% of the time
	t.Skip()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.cosmosAccount(data, true, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cosmosAccount(data, true, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cosmosAccount(data, false, false),
			Check: resource.ComposeTestCheckFunc(
				testCheckAdvancedThreatProtectionIsFalse(data.ResourceName),
			),
		},
	})
}

func TestAccAdvancedThreatProtection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advanced_threat_protection", "test")
	r := AdvancedThreatProtectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageAccount(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t AdvancedThreatProtectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AdvancedThreatProtectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.AdvancedThreatProtectionClient.Get(ctx, id.TargetResourceID)
	if err != nil {
		return nil, fmt.Errorf("reading Advanced Threat Protection (%s): %+v", id.TargetResourceID, err)
	}

	return utils.Bool(resp.AdvancedThreatProtectionProperties != nil), nil
}

// nolint unused
func testCheckAdvancedThreatProtectionIsFalse(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		targetResourceId := rs.Primary.Attributes["target_resource_id"]

		resp, err := client.Get(ctx, targetResourceId)
		if err != nil {
			return fmt.Errorf("Failed reading Advanced Threat Protection for resource %q: %+v", targetResourceId, err)
		}

		if props := resp.AdvancedThreatProtectionProperties; props != nil {
			if props.IsEnabled != nil {
				if *props.IsEnabled {
					return fmt.Errorf("Advanced Threat Protection is still true for resource %q: %+v", targetResourceId, err)
				}
			}
		}

		return nil
	}
}

func (AdvancedThreatProtectionResource) requiresImport(data acceptance.TestData) string {
	template := AdvancedThreatProtectionResource{}.storageAccount(data, true)
	return fmt.Sprintf(`
%s

resource "azurerm_advanced_threat_protection" "import" {
  target_resource_id = azurerm_advanced_threat_protection.test.target_resource_id
  enabled            = azurerm_advanced_threat_protection.test.enabled
}
`, template)
}

func (AdvancedThreatProtectionResource) storageAccount(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ATP-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_advanced_threat_protection" "test" {
  target_resource_id = "${azurerm_storage_account.test.id}"
  enabled            = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, enabled)
}

// nolint unused - mistakenly marked as unused
func (AdvancedThreatProtectionResource) cosmosAccount(data acceptance.TestData, hasResource, enabled bool) string {
	atp := ""
	if hasResource {
		atp = fmt.Sprintf(`
resource "azurerm_advanced_threat_protection" "test" {
  target_resource_id = "${azurerm_cosmosdb_account.test.id}"
  enabled            = %t
}
`, enabled)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ATP-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomString, atp)
}
