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

func resourceSecurityCenterSettingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).SettingClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
