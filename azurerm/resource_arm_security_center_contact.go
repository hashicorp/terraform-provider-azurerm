package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/2017-08-01-preview/security"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
)

//seems you can only set one contact:
// Invalid security contact name was provided - only 'defaultX' is allowed where X is an index
// Invalid security contact name 'default0' was provided. Expected 'default1'
// Message="Invalid security contact name 'default2' was provided. Expected 'default1'"
const securityCenterContactName = "default1"

func resourceArmSecurityCenterContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterContactCreateUpdate,
		Read:   resourceArmSecurityCenterContactRead,
		Update: resourceArmSecurityCenterContactCreateUpdate,
		Delete: resourceArmSecurityCenterContactDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"phone": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"alert_notifications": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alerts_to_admins": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceArmSecurityCenterContactCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenterContactsClient
	ctx := meta.(*ArmClient).StopContext

	name := securityCenterContactName

	contact := security.Contact{
		ContactProperties: &security.ContactProperties{
			Email: utils.String(d.Get("email").(string)),
			Phone: utils.String(d.Get("phone").(string)),
		},
	}

	if alertNotifications := d.Get("alert_notifications").(bool); alertNotifications {
		contact.AlertNotifications = security.On
	} else {
		contact.AlertNotifications = security.Off
	}

	if alertNotifications := d.Get("alerts_to_admins").(bool); alertNotifications {
		contact.AlertsToAdmins = security.AlertsToAdminsOn
	} else {
		contact.AlertsToAdmins = security.AlertsToAdminsOff
	}

	if d.IsNewResource() {
		_, err := client.Create(ctx, name, contact)
		if err != nil {
			return fmt.Errorf("Error creating Security Center Contact: %+v", err)
		}

		resp, err := client.Get(ctx, name)
		if err != nil {
			return fmt.Errorf("Error reading Security Center Contact: %+v", err)
		}
		if resp.ID == nil {
			return fmt.Errorf("Security Center Contact ID is nil")
		}

		d.SetId(*resp.ID)
	} else {
		_, err := client.Update(ctx, name, contact)
		if err != nil {
			return fmt.Errorf("Error updating Security Center Contact: %+v", err)
		}
	}

	return resourceArmSecurityCenterContactRead(d, meta)
}

func resourceArmSecurityCenterContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenterContactsClient
	ctx := meta.(*ArmClient).StopContext

	name := securityCenterContactName

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Security Center Contact: %+v", err)
	}

	if properties := resp.ContactProperties; properties != nil {
		d.Set("email", properties.Email)
		d.Set("phone", properties.Phone)
		d.Set("alert_notifications", properties.AlertNotifications == security.On)
		d.Set("alerts_to_admins", properties.AlertsToAdmins == security.AlertsToAdminsOn)
	}

	return nil
}

func resourceArmSecurityCenterContactDelete(_ *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenterContactsClient
	ctx := meta.(*ArmClient).StopContext

	name := securityCenterContactName

	resp, err := client.Delete(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Error deleting Security Center Contact: %+v", err)
	}

	return nil
}
