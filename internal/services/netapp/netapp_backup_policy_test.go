// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backuppolicy"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppBackupPolicyResource struct{}

func (t NetAppBackupPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := backuppolicy.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BackupPolicyClient.BackupPoliciesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccNetAppBackupPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_backup_policy", "test")
	r := NetAppBackupPolicyResource{}

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

func TestAccNetAppBackupPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_backup_policy", "test")
	r := NetAppBackupPolicyResource{}

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

func TestAccNetAppBackupPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_backup_policy", "test")
	r := NetAppBackupPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("2"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("2"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("2"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("10"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("10"),
				check.That(data.ResourceName).Key("daily_backups_to_keep").HasValue("10"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (r NetAppBackupPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_policy" "test" {
  name                = "acctest-NetAppBackupPolicy-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  account_name        = azurerm_netapp_account.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppBackupPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_policy" "test" {
  name                    = "acctest-NetAppBackupPolicy-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  account_name            = azurerm_netapp_account.test.name
  daily_backups_to_keep   = 2
  weekly_backups_to_keep  = 2
  monthly_backups_to_keep = 2
  enabled                 = true

  tags = {
    "testTag" = "testTagValue"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppBackupPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_policy" "test" {
  name                    = "acctest-NetAppBackupPolicy-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  account_name            = azurerm_netapp_account.test.name
  daily_backups_to_keep   = 10
  weekly_backups_to_keep  = 10
  monthly_backups_to_keep = 10
  enabled                 = false

  tags = {
    "testTag" = "testTagValue",
    "FoO"     = "BaR"
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetAppBackupPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}

`, data.RandomInteger, data.Locations.Primary)
}
