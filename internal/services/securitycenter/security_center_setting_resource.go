// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// TODO: this resource should be split into data_export_setting and alert_sync_setting

func resourceSecurityCenterSetting() *pluginsdk.Resource {
	validSettingName := settings.PossibleValuesForSettingName()

	if !features.FourPointOhBeta() {
		// This is for backward compatibility.. The swagger defines the valid enum to be "Sensinel" (see below), so this ("SENTINEL") shall be removed since 4.0.
		// https://github.com/Azure/azure-rest-api-specs/blob/b52464f520b77222ac8b0bdeb80a030c0fdf5b1b/specification/security/resource-manager/Microsoft.Security/stable/2021-06-01/settings.json#L285
		validSettingName = append(validSettingName, "SENTINEL")
	}

	return &pluginsdk.Resource{
		Create: resourceSecurityCenterSettingUpdate,
		Read:   resourceSecurityCenterSettingRead,
		Update: resourceSecurityCenterSettingUpdate,
		Delete: resourceSecurityCenterSettingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SettingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SecurityCenterSettingsV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"setting_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func() func(string, string, string, *schema.ResourceData) bool {
					// This is a workaround for `SENTINEL` value.
					if !features.FourPointOhBeta() {
						return suppress.CaseDifference
					}
					return nil
				}(),
				ValidateFunc: validation.StringInSlice(validSettingName, false),
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSecurityCenterSettingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.SettingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	settingName := d.Get("setting_name").(string)

	if !features.FourPointOhBeta() && settingName == "SENTINEL" {
		settingName = "Sentinel"
	}

	id := settings.NewSettingID(subscriptionId, settings.SettingName(settingName))

	if d.IsNewResource() {
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

func resourceSecurityCenterSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.SettingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := settings.ParseSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model != nil {
		if alertSyncSettings, ok := resp.Model.(settings.AlertSyncSettings); ok && alertSyncSettings.Properties != nil {
			d.Set("enabled", alertSyncSettings.Properties.Enabled)
		}
		if dataExportSettings, ok := resp.Model.(settings.DataExportSettings); ok && dataExportSettings.Properties != nil {
			d.Set("enabled", dataExportSettings.Properties.Enabled)
		}
	}

	d.Set("setting_name", id.SettingName)

	return nil
}

func resourceSecurityCenterSettingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.SettingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := settings.ParseSettingID(d.Id())
	if err != nil {
		return err
	}

	setting, err := expandSecurityCenterSetting(id.SettingName, false)
	if err != nil {
		return err
	}

	if _, err := client.Update(ctx, *id, setting); err != nil {
		return fmt.Errorf("disabling %s: %+v", id, err)
	}

	return nil
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
