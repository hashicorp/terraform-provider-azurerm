// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
		return pointer.To(false), nil
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
			Config: r.update(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.update(data, false),
			ExpectError: regexp.MustCompile("`malware_scanning_on_upload_filters` cannot be set if `malware_scanning_on_upload_enabled` is `false`"),
		},
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
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

func TestAccSecurityCenterStorageDefender_invalidFilterBlobPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_storage_defender", "test")
	r := SecurityCenterStorageDefenderResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidFilterBlobPrefix(data),
			ExpectError: regexp.MustCompile(regexp.QuoteMeta(`expected value of malware_scanning_on_upload_filters.0.exclude_blobs_with_prefix.0 to not contain any of "\\?#%&+:*\"|", got test/\`)),
		},
	})
}

func (r SecurityCenterStorageDefenderResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

func (r SecurityCenterStorageDefenderResource) update(data acceptance.TestData, malwareScanningOnUploadEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id                             = azurerm_storage_account.test.id
  override_subscription_settings_enabled         = false
  malware_scanning_on_upload_enabled             = %t
  malware_scanning_on_upload_cap_gb_per_month    = 6
  sensitive_data_discovery_enabled               = false
  malware_scanning_write_results_on_tags_enabled = false

  malware_scanning_on_upload_filters {
    exclude_blobs_larger_than = 131072

    exclude_blobs_with_prefix = [
      "container-2/blob",
      "container-3/blob-0"
    ]
  }
}
`, r.template(data), malwareScanningOnUploadEnabled)
}

func (r SecurityCenterStorageDefenderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id                             = azurerm_storage_account.test.id
  override_subscription_settings_enabled         = true
  malware_scanning_on_upload_enabled             = true
  malware_scanning_on_upload_cap_gb_per_month    = 4
  malware_scanning_write_results_on_tags_enabled = true
  sensitive_data_discovery_enabled               = true

  malware_scanning_on_upload_filters {
    exclude_blobs_larger_than = 65536

    exclude_blobs_with_prefix = [
      "container-0",
      "container-1/",
      "container-2/blob",
      "container-3/blob-0"
    ]

    exclude_blobs_with_suffix = [
      ".jpg",
      ".cpkt.index"
    ]
  }
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

func (r SecurityCenterStorageDefenderResource) invalidFilterBlobPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_storage_defender" "test" {
  storage_account_id = azurerm_storage_account.test.id

  malware_scanning_on_upload_filters {
    exclude_blobs_with_prefix = [
      "test/\\"
    ]
  }
}
`, r.template(data))
}
