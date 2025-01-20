// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsWorkspaceResource struct {
	DoNotRunFreeTierTest bool
}

func TestAccLogAnalyticsWorkspace_basic(t *testing.T) {
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

func TestAccLogAnalyticsWorkspace_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{DoNotRunFreeTierTest: true}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.freeTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.freeTierWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogAnalyticsWorkspaceResource) preCheck(t *testing.T) {
	if r.DoNotRunFreeTierTest {
		t.Skipf("`TestAccLogAnalyticsWorkspace_freeTier` due to subscription pricing tier configuration (e.g. Pricing tier doesn't match the subscription's billing model.)")
	}
}

func (LogAnalyticsWorkspaceResource) freeTierWithTags(data acceptance.TestData) string {
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
  sku                 = "Free"
  retention_in_days   = 7

  tags = {
    Environment = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
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
		{
			Config: r.withUseResourceOnlyPermission(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withUseResourceOnlyPermission(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLogAnalyticsWorkspace_ToggleDisableLocalAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDisableLocalAuth(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withDisableLocalAuth(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withDisableLocalAuth(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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
		{
			Config: r.withImmediatePurgeDataOn30Days(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withImmediatePurgeDataOn30Days(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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

	return utils.Bool(resp.Model != nil), nil
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

func (LogAnalyticsWorkspaceResource) freeTier(data acceptance.TestData) string {
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
  sku                 = "Free"
  retention_in_days   = 7
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

func (LogAnalyticsWorkspaceResource) withDisableLocalAuth(data acceptance.TestData, disableLocalAuth bool) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, disableLocalAuth)
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
