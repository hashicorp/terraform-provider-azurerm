package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceArmSecurityCenterSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterSettingUpdate,
		Read:   resourceArmSecurityCenterSettingRead,
		Update: resourceArmSecurityCenterSettingUpdate,
		Delete: resourceArmSecurityCenterSettingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"setting_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MCAS",
					"WDATP",
				}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceArmSecurityCenterSettingUpdate(d *schema.ResourceData, meta interface{}) error {
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

	return resourceArmSecurityCenterSettingRead(d, meta)
}

func resourceArmSecurityCenterSettingRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceArmSecurityCenterSettingDelete(_ *schema.ResourceData, _ interface{}) error {
	log.Printf("[DEBUG] Security Center deletion invocation")
	return nil // cannot be deleted.
}
