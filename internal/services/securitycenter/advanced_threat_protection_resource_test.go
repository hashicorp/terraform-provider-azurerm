// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdvancedThreatProtectionResource struct{}

func TestAccAdvancedThreatProtection_cosmosAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advanced_threat_protection", "test")
	r := AdvancedThreatProtectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosAccount(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cosmosAccount(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cosmosAccount(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(checkAdvancedThreatProtectionIsFalse, "azurerm_cosmosdb_account.test"),
			),
		},
	})
}

func TestAccAdvancedThreatProtection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advanced_threat_protection", "test")
	r := AdvancedThreatProtectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosAccount(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (AdvancedThreatProtectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AdvancedThreatProtectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.AdvancedThreatProtectionClient.Get(ctx, id.TargetResourceID)
	if err != nil {
		return nil, fmt.Errorf("reading Advanced Threat Protection (%s): %+v", id.TargetResourceID, err)
	}

	return pointer.To(resp.AdvancedThreatProtectionProperties != nil), nil
}

func checkAdvancedThreatProtectionIsFalse(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	targetResourceID := state.ID
	resp, err := clients.SecurityCenter.AdvancedThreatProtectionClient.Get(ctx, targetResourceID)
	if err != nil {
		return fmt.Errorf("reading Advanced Threat Protection (%s): %+v", targetResourceID, err)
	}

	if resp.AdvancedThreatProtectionProperties == nil || resp.IsEnabled == nil {
		return fmt.Errorf("Advanced Threat Protection (%s) properties is nil", targetResourceID)
	}

	if *resp.IsEnabled {
		return fmt.Errorf("Advanced Threat Protection (%s) properties is still enabled", targetResourceID)
	}

	return nil
}

func (AdvancedThreatProtectionResource) requiresImport(data acceptance.TestData) string {
	template := AdvancedThreatProtectionResource{}.cosmosAccount(data, true, true)
	return fmt.Sprintf(`
%s

resource "azurerm_advanced_threat_protection" "import" {
  target_resource_id = azurerm_advanced_threat_protection.test.target_resource_id
  enabled            = azurerm_advanced_threat_protection.test.enabled
}
`, template)
}

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
