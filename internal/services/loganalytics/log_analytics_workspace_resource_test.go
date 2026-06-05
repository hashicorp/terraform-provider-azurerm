// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceResource struct{}

func TestAccLogAnalyticsWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace"),
		},
	})
}

func TestAccLogAnalyticsWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

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

func TestAccLogAnalyticsWorkspace_withDefaultSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDefaultSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withVolumeCap(data, 4.5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_removeVolumeCap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withVolumeCap(data, 5.5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeVolumeCap(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("daily_quota_gb").HasValue("-1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withInternetIngestionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetIngestionEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetIngestionEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withInternetQueryEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetQueryEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetQueryEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_internetIngestionAccessType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetIngestionAccessType(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetIngestionAccessType(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetIngestionAccessType(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_internetQueryAccessType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withInternetQueryAccessType(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_query_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetQueryAccessType(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_query_access_type").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withInternetQueryAccessType(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_query_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_internetAccessTypeDeprecatedBool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	if features.FivePointOh() {
		t.Skip("Skipping since internet_ingestion_enabled and internet_query_enabled are removed in v5.0")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Deprecated bool = true: access_type must be computed as "Enabled"
			Config: r.withInternetIngestionEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
		{
			// Deprecated bool = false: access_type must be computed as "Disabled"
			Config: r.withInternetIngestionEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
		{
			// Migrate from deprecated bool to new string attr — no diff expected
			Config: r.withInternetIngestionAccessType(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_internetAccessTypeSecuredByPerimeter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Associate workspace with an NSP in Enforced mode — Azure sets both
			// public network access types to SecuredByPerimeter automatically.
			Config: r.withNspEnforced(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("SecuredByPerimeter"),
				check.That(data.ResourceName).Key("internet_query_access_type").HasValue("SecuredByPerimeter"),
			),
		},
		data.ImportStep(),
		{
			// Remove the NSP association and explicitly set both back to Enabled.
			Config: r.withNspRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("internet_ingestion_access_type").HasValue("Enabled"),
				check.That(data.ResourceName).Key("internet_query_access_type").HasValue("Enabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withCapacityReservation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCapacityReservation(data, 2000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("reservation_capacity_in_gb_per_day").HasValue("2000"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCapacityReservationTypo(data, 5000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("reservation_capacity_in_gb_per_day").HasValue("5000"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_negativeOne(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withVolumeCap(data, -1.0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_cmkForQueryForced(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkForQueryForced(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cmkForQueryForced(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_ToggleAllowOnlyResourcePermission(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUseResourceOnlyPermission(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUseResourceOnlyPermission(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUseResourceOnlyPermission(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_ToggleEnableLocalAuthDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	if features.FivePointOh() {
		t.Skip("Skipping since local_authentication_disabled is deprecated in 5.0")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.localAuthEnabledDeprecated(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.localAuthEnabledDeprecated(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("true"),
			),
		},
		data.ImportStep(),

		{
			Config: r.localAuthEnabledDeprecated(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_authentication_disabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_ToggleEnableLocalAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.localAuthEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.localAuthEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.localAuthEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("local_authentication_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_ToggleImmediatelyPurgeDataOn30Days(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withImmediatePurgeDataOn30Days(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withImmediatePurgeDataOn30Days(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withImmediatePurgeDataOn30Days(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCapacityReservationTypo(data, 100),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("reservation_capacity_in_gb_per_day").HasValue("100"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withDefaultDataCollectionRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	// the default data collection rule could only be set during update,
	// and to avoid the dependency cycle, we do an additional update here.
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDataCollectionRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withDefaultDataCollectionRule(data, datacollectionrules.NewDataCollectionRuleID(data.Subscriptions.Primary, fmt.Sprintf("acctestRG-%d", data.RandomInteger), fmt.Sprintf("acctestmdcr-%d", data.RandomInteger)).ID()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withSystemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_withUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspace_toggleIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.SharedKeyWorkspacesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Workspace (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LogAnalyticsWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  identity {
    type = "SystemAssigned"
  }

}
`, data.RandomInteger, data.Locations.Primary)
}

func (LogAnalyticsWorkspaceResource) withUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctest-%[1]d"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogAnalyticsWorkspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "import" {
  name                = azurerm_log_analytics_workspace.test.name
  location            = azurerm_log_analytics_workspace.test.location
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
  sku                 = "PerGB2018"
}
`, r.basic(data))
}

func (LogAnalyticsWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withDefaultSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withVolumeCap(data acceptance.TestData, volumeCapGb float64) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
  daily_quota_gb      = %f

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, volumeCapGb)
}

func (LogAnalyticsWorkspaceResource) removeVolumeCap(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetIngestionEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                       = "acctestLAW-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  internet_ingestion_enabled = true
  sku                        = "PerGB2018"
  retention_in_days          = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetIngestionEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                       = "acctestLAW-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  internet_ingestion_enabled = false
  sku                        = "PerGB2018"
  retention_in_days          = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetQueryEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                   = "acctestLAW-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  internet_query_enabled = true
  sku                    = "PerGB2018"
  retention_in_days      = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withInternetQueryEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                   = "acctestLAW-%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  internet_query_enabled = false
  sku                    = "PerGB2018"
  retention_in_days      = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withCapacityReservation(data acceptance.TestData, capacityReservation int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                               = "acctestLAW-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  internet_query_enabled             = false
  sku                                = "CapacityReservation"
  reservation_capacity_in_gb_per_day = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacityReservation)
}

func (LogAnalyticsWorkspaceResource) withCapacityReservationTypo(data acceptance.TestData, capacityReservation int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                               = "acctestLAW-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  internet_query_enabled             = false
  sku                                = "CapacityReservation"
  reservation_capacity_in_gb_per_day = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacityReservation)
}

func (LogAnalyticsWorkspaceResource) cmkForQueryForced(data acceptance.TestData, cmkForQueryForced bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                 = "acctestLAW-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku                  = "PerGB2018"
  retention_in_days    = 30
  cmk_for_query_forced = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, cmkForQueryForced)
}

func (LogAnalyticsWorkspaceResource) withUseResourceOnlyPermission(data acceptance.TestData, useResourceOnlyPermission bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                            = "acctestLAW-%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  sku                             = "PerGB2018"
  retention_in_days               = 30
  allow_resource_only_permissions = %[4]t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, useResourceOnlyPermission)
}

func (LogAnalyticsWorkspaceResource) localAuthEnabledDeprecated(data acceptance.TestData, localAuthEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                          = "acctestLAW-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  sku                           = "PerGB2018"
  retention_in_days             = 30
  local_authentication_disabled = %[4]t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, !localAuthEnabled)
}

func (LogAnalyticsWorkspaceResource) localAuthEnabled(data acceptance.TestData, localAuthEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                         = "acctestLAW-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  sku                          = "PerGB2018"
  retention_in_days            = 30
  local_authentication_enabled = %[4]t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, localAuthEnabled)
}

func (LogAnalyticsWorkspaceResource) withDataCollectionRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "WorkspaceTransforms"

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
      name                  = "test-destination"
    }
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination"]
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LogAnalyticsWorkspaceResource) withDefaultDataCollectionRule(data acceptance.TestData, ruleID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%s"
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctestmdcr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  kind                = "WorkspaceTransforms"

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
      name                  = "test-destination"
    }
  }

  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination"]
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                    = "acctestLAW-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku                     = "PerGB2018"
  retention_in_days       = 30
  data_collection_rule_id = "%[3]s"
}
`, data.RandomInteger, data.Locations.Primary, ruleID)
}

func (LogAnalyticsWorkspaceResource) withImmediatePurgeDataOn30Days(data acceptance.TestData, enable bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                                    = "acctestLAW-%d"
  location                                = azurerm_resource_group.test.location
  resource_group_name                     = azurerm_resource_group.test.name
  sku                                     = "PerGB2018"
  retention_in_days                       = 30
  immediate_data_purge_on_30_days_enabled = %[4]t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enable)
}

func (LogAnalyticsWorkspaceResource) withInternetIngestionAccessType(data acceptance.TestData, accessType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                           = "acctestLAW-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  internet_ingestion_access_type = "%s"
  sku                            = "PerGB2018"
  retention_in_days              = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, accessType)
}

func (LogAnalyticsWorkspaceResource) withInternetQueryAccessType(data acceptance.TestData, accessType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                       = "acctestLAW-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  internet_query_access_type = "%s"
  sku                        = "PerGB2018"
  retention_in_days          = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, accessType)
}

func (LogAnalyticsWorkspaceResource) withNspEnforced(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name                = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_network_security_perimeter_profile" "test" {
  name                          = "acctestNspProfile-%d"
  network_security_perimeter_id = azurerm_network_security_perimeter.test.id
}

# The workspace is declared with SecuredByPerimeter so that Terraform state
# and the Azure API value agree once the NSP association below is in place.
resource "azurerm_log_analytics_workspace" "test" {
  name                           = "acctestLAW-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  sku                            = "PerGB2018"
  retention_in_days              = 30
  internet_ingestion_access_type = "SecuredByPerimeter"
  internet_query_access_type     = "SecuredByPerimeter"
}

resource "azurerm_network_security_perimeter_association" "test" {
  name                                  = "acctestNspAssoc-%d"
  network_security_perimeter_profile_id = azurerm_network_security_perimeter_profile.test.id
  resource_id                           = azurerm_log_analytics_workspace.test.id
  access_mode                           = "Enforced"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (LogAnalyticsWorkspaceResource) withNspRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                           = "acctestLAW-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  sku                            = "PerGB2018"
  retention_in_days              = 30
  internet_ingestion_access_type = "Enabled"
  internet_query_access_type     = "Enabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
