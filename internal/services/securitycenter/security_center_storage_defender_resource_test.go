// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterStorageDefenderResource struct{}

func (SecurityCenterStorageDefenderResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseScopeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.DefenderForStorageClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.IsEnabled == nil {
		return utils.Bool(false), nil
	}

	return resp.Model.Properties.IsEnabled, nil
}

func TestAccSecurityCenterStorageDefender_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterStorageDefender_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterStorageDefender_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterStorageDefender_reapply(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSecurityCenterStorageDefender_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSecurityCenterStorageDefender_eventGrid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventGrid(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SecurityCenterStorageDefenderResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestacc%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r SecurityCenterStorageDefenderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id = azurerm_storage_account.test.id
}
`, r.template(data))
}

func (r SecurityCenterStorageDefenderResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id                          = azurerm_storage_account.test.id
  override_subscription_settings_enabled      = false
  malware_scanning_on_upload_enabled          = false
  malware_scanning_on_upload_cap_gb_per_month = 6
  sensitive_data_discovery_enabled            = false
}
`, r.template(data))
}

func (r SecurityCenterStorageDefenderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id                          = azurerm_storage_account.test.id
  override_subscription_settings_enabled      = true
  malware_scanning_on_upload_enabled          = true
  malware_scanning_on_upload_cap_gb_per_month = 4
  sensitive_data_discovery_enabled            = true
}
`, r.template(data))
}

func (r SecurityCenterStorageDefenderResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "import" {
  storage_account_id = azurerm_security_center_storage_defender.test.id
}
`, r.basic(data))
}

func (r SecurityCenterStorageDefenderResource) eventGrid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctestEVGT-storage-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id                     = azurerm_storage_account.test.id
  override_subscription_settings_enabled = true
  malware_scanning_on_upload_enabled     = true
  scan_results_event_grid_topic_id       = azurerm_eventgrid_topic.test.id
}
`, r.template(data), data.RandomInteger)
}
