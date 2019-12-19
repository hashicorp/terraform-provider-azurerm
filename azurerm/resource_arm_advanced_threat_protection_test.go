package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
)

func TestAccAzureRMAdvancedThreatProtection_storageAccount(t *testing.T) {
	rn := "azurerm_advanced_threat_protection.test"
	ri := tf.AccRandTimeInt()
	var id securitycenter.AdvancedThreatProtectionResourceID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvancedThreatProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvancedThreatProtection_storageAccount(ri, acceptance.Location(), true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionExists(rn, &id),
					resource.TestCheckResourceAttr(rn, "enabled", "true"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMAdvancedThreatProtection_storageAccount(ri, acceptance.Location(), true, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionExists(rn, &id),
					resource.TestCheckResourceAttr(rn, "enabled", "false"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMAdvancedThreatProtection_storageAccount(ri, acceptance.Location(), false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionIsFalse(&id),
				),
			},
		},
	})
}

func TestAccAzureRMAdvancedThreatProtection_cosmosAccount(t *testing.T) {
	rn := "azurerm_advanced_threat_protection.test"
	ri := tf.AccRandTimeInt()
	var id securitycenter.AdvancedThreatProtectionResourceID

	// the API errors on deleting the cosmos DB account some of the time so lets skip this test for now
	// TODO: remove once this is fixed: https://github.com/Azure/azure-sdk-for-go/issues/6310
	// run it multiple times in a row as it only fails 50% of the time
	t.Skip()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvancedThreatProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvancedThreatProtection_cosmosAccount(ri, acceptance.Location(), true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionExists(rn, &id),
					resource.TestCheckResourceAttr(rn, "enabled", "true"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMAdvancedThreatProtection_cosmosAccount(ri, acceptance.Location(), true, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionExists(rn, &id),
					resource.TestCheckResourceAttr(rn, "enabled", "false"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMAdvancedThreatProtection_cosmosAccount(ri, acceptance.Location(), false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionIsFalse(&id),
				),
			},
		},
	})
}

func TestAccAzureRMAdvancedThreatProtection_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	rn := "azurerm_advanced_threat_protection.test"
	ri := tf.AccRandTimeInt()
	var id securitycenter.AdvancedThreatProtectionResourceID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAdvancedThreatProtectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdvancedThreatProtection_storageAccount(ri, acceptance.Location(), true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdvancedThreatProtectionExists(rn, &id),
					resource.TestCheckResourceAttr(rn, "enabled", "true"),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccAzureRMAdvancedThreatProtection_requiresImport(ri, acceptance.Location()),
				ExpectError: acceptance.RequiresImportError("azurerm_advanced_threat_protection"),
			},
		},
	})
}

func testCheckAzureRMAdvancedThreatProtectionExists(resourceName string, idToReturn *securitycenter.AdvancedThreatProtectionResourceID) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure resource group exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := securitycenter.ParseAdvancedThreatProtectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.TargetResourceID)
		if err != nil {
			return fmt.Errorf("Bad: Get on AdvancedThreatProtectionClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Advanced Threat Protection for resource %q not found", id.TargetResourceID)
		}

		*idToReturn = *id

		return nil
	}
}

func testCheckAzureRMAdvancedThreatProtectionIsFalse(id *securitycenter.AdvancedThreatProtectionResourceID) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, id.TargetResourceID)
		if err != nil {
			return fmt.Errorf("Failed reading Advanced Threat Protection for resource %q: %+v", id.TargetResourceID, err)
		}

		if props := resp.AdvancedThreatProtectionProperties; props != nil {
			if props.IsEnabled != nil {
				if *props.IsEnabled {
					return fmt.Errorf("Advanced Threat Protection is still true for resource %q: %+v", id.TargetResourceID, err)
				}
			}
		}

		return nil
	}
}

func testCheckAzureRMAdvancedThreatProtectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_advanced_threat_protection" {
			continue
		}

		id, err := securitycenter.ParseAdvancedThreatProtectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.TargetResourceID)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Advanced Threat Protection still exists:\n%#v", resp.ID)
		}
	}

	return nil
}

func testAccAzureRMAdvancedThreatProtection_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_advanced_threat_protection" "requireimport" {
  target_resource_id = "${azurerm_advanced_threat_protection.test.target_resource_id}"
  enabled            = "${azurerm_advanced_threat_protection.test.enabled}"
}
`, testAccAzureRMAdvancedThreatProtection_storageAccount(rInt, location, true, true))
}

func testAccAzureRMAdvancedThreatProtection_storageAccount(rInt int, location string, hasResource, enabled bool) string {
	atp := ""
	if hasResource {
		atp = fmt.Sprintf(`
resource "azurerm_advanced_threat_protection" "test" {
  target_resource_id = "${azurerm_storage_account.test.id}"
  enabled            = %t
}
`, enabled)
	}

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ATP-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctest%[3]d"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

%[4]s
`, rInt, location, rInt/10, atp)
}

func testAccAzureRMAdvancedThreatProtection_cosmosAccount(rInt int, location string, hasResource, enabled bool) string {
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
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ATP-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "%[3]s"
  }

  geo_location {
    location          = "${azurerm_resource_group.test.location}"
    failover_priority = 0
  }
}

%[4]s
`, rInt, location, string(documentdb.Eventual), atp)
}
