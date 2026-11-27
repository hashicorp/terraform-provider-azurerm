// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_setting

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSettingUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.SettingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	settingName := d.Get("setting_name").(string)
	id := settings.NewSettingID(subscriptionId, settings.SettingName(settingName))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %v", id, err)
			}

			if existing.Model != nil {
				if alertSyncSettings, ok := existing.Model.(settings.AlertSyncSettings); ok && alertSyncSettings.Properties != nil && alertSyncSettings.Properties.Enabled {
					return tf.ImportAsExistsError("azurerm_security_center_setting", id.ID())
				}
				if dataExportSettings, ok := existing.Model.(settings.DataExportSettings); ok && dataExportSettings.Properties != nil && dataExportSettings.Properties.Enabled {
					return tf.ImportAsExistsError("azurerm_security_center_setting", id.ID())
				}
			}
		}
	}

	setting, err := expandSecurityCenterSetting(id.SettingName, d.Get("enabled").(bool))
	if err != nil {
		return err
	}

	if _, err := client.Update(ctx, id, setting); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSettingRead(d, meta)
}

func expandSecurityCenterSetting(name settings.SettingName, enabled bool) (settings.Setting, error) {
	switch name {
	case settings.SettingNameMCAS,
		settings.SettingNameWDATP,
		settings.SettingNameWDATPEXCLUDELINUXPUBLICPREVIEW,
		settings.SettingNameWDATPUNIFIEDSOLUTION:
		return settings.DataExportSettings{
			Properties: &settings.DataExportSettingProperties{
				Enabled: enabled,
			},
		}, nil
	case "SENTINEL",
		settings.SettingNameSentinel:
		return settings.AlertSyncSettings{
			Properties: &settings.AlertSyncSettingProperties{
				Enabled: enabled,
			},
		}, nil
	default:
		return nil, fmt.Errorf("failed to deduce the kind from its name %q", name)
	}
}
