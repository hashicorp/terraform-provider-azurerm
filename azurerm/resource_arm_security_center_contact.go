package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := securityCenterContactName

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Security Center Contact: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_contact", *existing.ID)
		}
	}

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
		if _, err := client.Create(ctx, name, contact); err != nil {
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
		if _, err := client.Update(ctx, name, contact); err != nil {
			return fmt.Errorf("Error updating Security Center Contact: %+v", err)
		}
	}

	return resourceArmSecurityCenterContactRead(d, meta)
}

func resourceArmSecurityCenterContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

func resourceArmSecurityCenterContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
