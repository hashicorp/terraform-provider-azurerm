// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_setting

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-05-01/settings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSettingRead(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).SettingClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
