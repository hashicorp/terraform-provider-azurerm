// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// TODO: this resource should be split into data_export_setting and alert_sync_setting

func resourceSecurityCenterSetting() *pluginsdk.Resource {
	validSettingName := settings.PossibleValuesForSettingName()

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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
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
	id := settings.NewSettingID(subscriptionId, settings.SettingName(settingName))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %v", id, err)
		}

		if existing.Model != nil {
			if alertSyncSettings, ok := existing.Model.(settings.AlertSyncSettings); ok {
				properties, err := expandAlertSyncSettingProperties(alertSyncSettings.Properties)
				if err != nil {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if properties != nil && properties.Enabled {
					return tf.ImportAsExistsError("azurerm_security_center_setting", id.ID())
				}
			}
			if dataExportSettings, ok := existing.Model.(settings.DataExportSettings); ok {
				properties, err := expandDataExportSettingProperties(dataExportSettings.Properties)
				if err != nil {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if properties != nil && properties.Enabled {
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
		if alertSyncSettings, ok := resp.Model.(settings.AlertSyncSettings); ok {
			properties, err := expandAlertSyncSettingProperties(alertSyncSettings.Properties)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if properties != nil {
				d.Set("enabled", properties.Enabled)
			}
		}
		if dataExportSettings, ok := resp.Model.(settings.DataExportSettings); ok {
			properties, err := expandDataExportSettingProperties(dataExportSettings.Properties)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if properties != nil {
				d.Set("enabled", properties.Enabled)
			}
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
			Properties: pointer.To(interface{}(settings.DataExportSettingProperties{
				Enabled: enabled,
			})),
		}, nil
	case "SENTINEL",
		settings.SettingNameSentinel:
		return settings.AlertSyncSettings{
			Properties: pointer.To(interface{}(settings.AlertSyncSettingProperties{
				Enabled: enabled,
			})),
		}, nil
	default:
		return nil, fmt.Errorf("failed to deduce the kind from its name %q", name)
	}
}

func expandAlertSyncSettingProperties(input *interface{}) (*settings.AlertSyncSettingProperties, error) {
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

func expandDataExportSettingProperties(input *interface{}) (*settings.DataExportSettingProperties, error) {
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
