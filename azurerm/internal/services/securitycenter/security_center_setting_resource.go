package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceSecurityCenterSetting() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterSettingUpdate,
		Read:   resourceSecurityCenterSettingRead,
		Update: resourceSecurityCenterSettingUpdate,
		Delete: resourceSecurityCenterSettingDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"setting_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MCAS",
					"WDATP",
				}, false),
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
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	settingName := d.Get("setting_name").(string)
	enabled := d.Get("enabled").(bool)
	setting := security.DataExportSettings{
		DataExportSettingProperties: &security.DataExportSettingProperties{
			Enabled: &enabled,
		},
		Kind: security.KindDataExportSettings,
	}

	if _, err := client.Update(ctx, settingName, setting); err != nil {
		return fmt.Errorf("Creating/updating Security Center pricing: %+v", err)
	}
	// TODO: switch to back when Swagger/API bug has been fixed:
	// https://github.com/Azure/azure-sdk-for-go/issues/12724
	resp, err := azuresdkhacks.GetSecurityCenterSetting(client, ctx, settingName)
	if err != nil {
		return fmt.Errorf("Reading Security Center setting: %+v", err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Nil/empty ID returned for Security Center setting %q", settingName)
	}

	d.SetId(*resp.ID)

	return resourceSecurityCenterSettingRead(d, meta)
}

func resourceSecurityCenterSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.SettingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SecurityCenterSettingID(d.Id())
	if err != nil {
		return err
	}
	// TODO: switch to back when Swagger/API bug has been fixed:
	// https://github.com/Azure/azure-sdk-for-go/issues/12724 (`Enabled` field missing)
	resp, err := azuresdkhacks.GetSecurityCenterSetting(client, ctx, id.SettingName)
	if err != nil {
		return fmt.Errorf("Reading Security Center setting: %+v", err)
	}

	if err != nil {
		return fmt.Errorf("Reading Security Center setting: %+v", err)
	}

	if properties := resp.DataExportSettingProperties; properties != nil {
		d.Set("enabled", properties.Enabled)
	}
	d.Set("setting_name", id.SettingName)

	return nil
}

func resourceSecurityCenterSettingDelete(_ *pluginsdk.ResourceData, _ interface{}) error {
	log.Printf("[DEBUG] Security Center deletion invocation")
	return nil // cannot be deleted.
}
