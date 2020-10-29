package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: 'default' is the only valid name currently supported by the API
// No other names can be created and the 'default' resource can not be destroyed
const securityCenterAutoProvisioningName = "default"

func resourceArmSecurityCenterAutoProvisioning() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterAutoProvisioningUpdate,
		Read:   resourceArmSecurityCenterAutoProvisioningRead,
		Update: resourceArmSecurityCenterAutoProvisioningUpdate,
		Delete: resourceArmSecurityCenterAutoProvisioningDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"auto_provision": {
				Type:     schema.TypeString,
				Required: true,
				// NOTE: the API seems case insensitive to this string value, 'ON', 'On', 'on' all work
				ValidateFunc: validation.StringInSlice([]string{
					string(security.AutoProvisionOn),
					string(security.AutoProvisionOff),
				}, false),
			},
		},
	}
}

func resourceArmSecurityCenterAutoProvisioningUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// No need for import check as there's always single resource called 'default'
	// - it cannot be deleted, all this does is set a string property to: "on" or "off"

	// Build settings struct with auto_provision value
	settings := security.AutoProvisioningSetting{
		AutoProvisioningSettingProperties: &security.AutoProvisioningSettingProperties{
			AutoProvision: security.AutoProvision(d.Get("auto_provision").(string)),
		},
	}

	// There is no update function or operation in the API, only create
	if _, err := client.Create(ctx, securityCenterAutoProvisioningName, settings); err != nil {
		return fmt.Errorf("Error creating/updating Security Center auto provisioning: %+v", err)
	}

	resp, err := client.Get(ctx, securityCenterAutoProvisioningName)
	if err != nil {
		return fmt.Errorf("Error reading Security Center auto provisioning: %+v", err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Security Center auto provisioning ID is nil or empty")
	}

	d.SetId(*resp.ID)

	return resourceArmSecurityCenterAutoProvisioningRead(d, meta)
}

func resourceArmSecurityCenterAutoProvisioningRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, securityCenterAutoProvisioningName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center subscription was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Security Center auto provisioning: %+v", err)
	}

	if properties := resp.AutoProvisioningSettingProperties; properties != nil {
		d.Set("auto_provision", properties.AutoProvision)
	}

	return nil
}

func resourceArmSecurityCenterAutoProvisioningDelete(d *schema.ResourceData, meta interface{}) error {
	// The API has no delete operation
	// Instead we reset back to 'Off' which is the default

	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	settings := security.AutoProvisioningSetting{
		AutoProvisioningSettingProperties: &security.AutoProvisioningSettingProperties{
			AutoProvision: "Off",
		},
	}

	// There is no update function or operation in the API, only create
	if _, err := client.Create(ctx, securityCenterAutoProvisioningName, settings); err != nil {
		return fmt.Errorf("Error resetting Security Center auto provisioning to 'Off': %+v", err)
	}

	return nil
}
