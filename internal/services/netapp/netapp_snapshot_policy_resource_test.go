// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/snapshotpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppSnapshotPolicyResource struct{}

func TestAccNetAppSnapshotPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

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

func TestAccNetAppSnapshotPolicy_hourlySchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hourlySchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshotPolicy_dailySchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dailySchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshotPolicy_weeklySchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.weeklySchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshotPolicy_monthlySchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.monthlySchedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshotPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_netapp_snapshot_policy"),
		},
	})
}

func TestAccNetAppSnapshotPolicy_updateSnapshotPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hourly_schedule.0.snapshots_to_keep").HasValue("1"),
				check.That(data.ResourceName).Key("daily_schedule.0.hour").HasValue("22"),
				check.That(data.ResourceName).Key("weekly_schedule.0.days_of_week.#").HasValue("2"),
				check.That(data.ResourceName).Key("monthly_schedule.0.days_of_month.#").HasValue("3"),
				check.That(data.ResourceName).Key("monthly_schedule.0.minute").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hourly_schedule.0.snapshots_to_keep").HasValue("5"),
				check.That(data.ResourceName).Key("daily_schedule.0.hour").HasValue("20"),
				check.That(data.ResourceName).Key("weekly_schedule.0.days_of_week.#").HasValue("3"),
				check.That(data.ResourceName).Key("monthly_schedule.0.days_of_month.#").HasValue("2"),
				check.That(data.ResourceName).Key("monthly_schedule.0.minute").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppSnapshotPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyResource{}

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

func (t NetAppSnapshotPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := snapshotpolicy.ParseSnapshotPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.SnapshotPoliciesClient.SnapshotPoliciesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp SnapshotPolicy (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetAppSnapshotPolicyResource) basic(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  hourly_schedule {
    snapshots_to_keep = 1
    minute            = 15
  }

  daily_schedule {
    snapshots_to_keep = 1
    hour              = 22
    minute            = 15
  }

  weekly_schedule {
    snapshots_to_keep = 1
    days_of_week      = ["Monday", "Friday"]
    hour              = 23
    minute            = 0
  }

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [1, 15, 30]
    hour              = 5
    minute            = 0
  }
}
`, template, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) hourlySchedule(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  hourly_schedule {
    snapshots_to_keep = 1
    minute            = 15
  }
}
`, template, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) dailySchedule(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  daily_schedule {
    snapshots_to_keep = 1
    hour              = 22
    minute            = 15
  }
}
`, template, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) weeklySchedule(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  weekly_schedule {
    snapshots_to_keep = 1
    days_of_week      = ["Monday", "Friday"]
    hour              = 23
    minute            = 0
  }
}
`, template, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) monthlySchedule(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [1, 15, 30]
    hour              = 5
    minute            = 0
  }
}
`, template, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) update(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  hourly_schedule {
    snapshots_to_keep = 5
    minute            = 15
  }

  daily_schedule {
    snapshots_to_keep = 1
    hour              = 20
    minute            = 15
  }

  weekly_schedule {
    snapshots_to_keep = 1
    days_of_week      = ["Monday", "Tuesday", "Friday"]
    hour              = 23
    minute            = 0
  }

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [1, 30]
    hour              = 5
    minute            = 30
  }
}
`, template, data.RandomInteger)
}

func (r NetAppSnapshotPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "import" {
  name                = azurerm_netapp_snapshot_policy.test.name
  location            = azurerm_netapp_snapshot_policy.test.location
  resource_group_name = azurerm_netapp_snapshot_policy.test.resource_group_name
  account_name        = azurerm_netapp_snapshot_policy.test.account_name
  enabled             = azurerm_netapp_snapshot_policy.test.enabled

  hourly_schedule {
    snapshots_to_keep = azurerm_netapp_snapshot_policy.test.hourly_schedule[0].snapshots_to_keep
    minute            = azurerm_netapp_snapshot_policy.test.hourly_schedule[0].minute
  }

  daily_schedule {
    snapshots_to_keep = azurerm_netapp_snapshot_policy.test.daily_schedule[0].snapshots_to_keep
    hour              = azurerm_netapp_snapshot_policy.test.daily_schedule[0].hour
    minute            = azurerm_netapp_snapshot_policy.test.daily_schedule[0].minute
  }

  weekly_schedule {
    snapshots_to_keep = azurerm_netapp_snapshot_policy.test.weekly_schedule[0].snapshots_to_keep
    days_of_week      = azurerm_netapp_snapshot_policy.test.weekly_schedule[0].days_of_week
    hour              = azurerm_netapp_snapshot_policy.test.weekly_schedule[0].hour
    minute            = azurerm_netapp_snapshot_policy.test.weekly_schedule[0].minute
  }

  monthly_schedule {
    snapshots_to_keep = azurerm_netapp_snapshot_policy.test.monthly_schedule[0].snapshots_to_keep
    days_of_month     = azurerm_netapp_snapshot_policy.test.monthly_schedule[0].days_of_month
    hour              = azurerm_netapp_snapshot_policy.test.monthly_schedule[0].hour
    minute            = azurerm_netapp_snapshot_policy.test.monthly_schedule[0].minute
  }
}
`, r.basic(data))
}

func (NetAppSnapshotPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Ternary, data.RandomInteger)
}

func (NetAppSnapshotPolicyResource) complete(data acceptance.TestData) string {
	template := NetAppSnapshotPolicyResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  hourly_schedule {
    snapshots_to_keep = 1
    minute            = 15
  }

  daily_schedule {
    snapshots_to_keep = 1
    hour              = 22
    minute            = 15
  }

  weekly_schedule {
    snapshots_to_keep = 1
    days_of_week      = ["Monday", "Friday"]
    hour              = 23
    minute            = 0
  }

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [1, 15, 30]
    hour              = 5
    minute            = 0
  }
  tags = {
    environment = "test"
  }
}
`, template, data.RandomInteger)
}
