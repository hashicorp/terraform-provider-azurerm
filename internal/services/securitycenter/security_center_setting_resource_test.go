// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SecurityCenterSettingResource struct{}

func TestAccSecurityCenterSetting(t *testing.T) {
	// there is only one workspace with the same name could exist, so run the tests in sequence.
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"setting": {
			"update":         testAccSecurityCenterSetting_update,
			"requiresImport": testAccSecurityCenterSetting_requiresImport,
		},
	})
}

func testAccSecurityCenterSetting_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_setting", "test")
	r := SecurityCenterSettingResource{}

	// lintignore:AT001
	testcases := []acceptance.TestStep{
		{
			Config: r.cfg("MCAS", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("MCAS"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("MCAS", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("setting_name").HasValue("MCAS"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("WDATP", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("setting_name").HasValue("WDATP"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("WDATP", false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("setting_name").HasValue("WDATP"),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.cfg("Sentinel", true),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
		{
			Config: r.cfg("Sentinel", false),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	}

	data.ResourceSequentialTest(t, r, testcases)
}

func testAccSecurityCenterSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_setting", "test")
	r := SecurityCenterSettingResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.cfg("MCAS", true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
		// reset
		{
			Config: r.cfg("MCAS", false),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (SecurityCenterSettingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := settings.ParseSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.SettingClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("checking for presence of existing %s: %v", id, err)
	}

	if resp.Model == nil {
		return pointer.To(false), nil
	}

	if alertSyncSettings, ok := resp.Model.(settings.AlertSyncSettings); ok {
		properties, err := expandAlertSyncSettingsTestProperties(alertSyncSettings.Properties)
		if err != nil {
			return nil, err
		}
		if properties == nil {
			return pointer.To(false), nil
		}
		return pointer.To(properties.Enabled), nil
	}
	if dataExportSettings, ok := resp.Model.(settings.DataExportSettings); ok {
		properties, err := expandDataExportSettingsTestProperties(dataExportSettings.Properties)
		if err != nil {
			return nil, err
		}
		if properties == nil {
			return pointer.To(false), nil
		}
		return pointer.To(properties.Enabled), nil
	}

	return pointer.To(false), nil
}

func (SecurityCenterSettingResource) Destroy(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.SecurityCenter.SettingClient
	id, err := settings.ParseSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	setting := settings.DataExportSettings{
		Properties: pointer.To(interface{}(settings.DataExportSettingProperties{
			Enabled: false,
		})),
	}

	if _, err := client.Update(ctx, *id, setting); err != nil {
		return nil, fmt.Errorf("disabling %s: %+v", id, err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("checking for presence of existing %s: %v", id, err)
	}

	if resp.Model == nil {
		return pointer.To(false), nil
	}

	if alertSyncSettings, ok := resp.Model.(settings.AlertSyncSettings); ok {
		properties, err := expandAlertSyncSettingsTestProperties(alertSyncSettings.Properties)
		if err != nil {
			return nil, err
		}
		if properties == nil || !properties.Enabled {
			return pointer.To(false), nil
		}
	}
	if dataExportSettings, ok := resp.Model.(settings.DataExportSettings); ok {
		properties, err := expandDataExportSettingsTestProperties(dataExportSettings.Properties)
		if err != nil {
			return nil, err
		}
		if properties == nil || !properties.Enabled {
			return pointer.To(false), nil
		}
	}

	return pointer.To(true), nil
}

func (SecurityCenterSettingResource) cfg(settingName string, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_setting" "test" {
  setting_name = "%s"
  enabled      = "%t"
}
`, settingName, enabled)
}

func (r SecurityCenterSettingResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_setting" "import" {
  setting_name = azurerm_security_center_setting.test.setting_name
  enabled      = azurerm_security_center_setting.test.enabled
}
`, r.cfg("MCAS", true))
}

func expandAlertSyncSettingsTestProperties(input *interface{}) (*settings.AlertSyncSettingProperties, error) {
	if input == nil || *input == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(*input)
	if err != nil {
		return nil, fmt.Errorf("marshaling alert sync setting properties: %+v", err)
	}

	var properties settings.AlertSyncSettingProperties
	if err := json.Unmarshal(encoded, &properties); err != nil {
		return nil, fmt.Errorf("unmarshaling alert sync setting properties: %+v", err)
	}

	return &properties, nil
}

func expandDataExportSettingsTestProperties(input *interface{}) (*settings.DataExportSettingProperties, error) {
	if input == nil || *input == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(*input)
	if err != nil {
		return nil, fmt.Errorf("marshaling data export setting properties: %+v", err)
	}

	var properties settings.DataExportSettingProperties
	if err := json.Unmarshal(encoded, &properties); err != nil {
		return nil, fmt.Errorf("unmarshaling data export setting properties: %+v", err)
	}

	return &properties, nil
}
