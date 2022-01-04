package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// seems you can only set one contact:
// Invalid security contact name was provided - only 'defaultX' is allowed where X is an index
// Invalid security contact name 'default0' was provided. Expected 'default1'
// Message="Invalid security contact name 'default2' was provided. Expected 'default1'"
const securityCenterContactName = "default1"

func resourceSecurityCenterContact() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterContactCreateUpdate,
		Read:   resourceSecurityCenterContactRead,
		Update: resourceSecurityCenterContactCreateUpdate,
		Delete: resourceSecurityCenterContactDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ContactID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"email": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"phone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alert_notifications": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"alerts_to_admins": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSecurityCenterContactCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	// TODO: split this Create/Update
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewContactID(subscriptionId, securityCenterContactName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.SubscriptionId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_security_center_contact", id.ID())
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
		// TODO: switch back when the Swagger/API bug has been fixed:
		// https://github.com/Azure/azure-rest-api-specs/issues/10717 (an undefined 201)
		if _, err := azuresdkhacks.CreateSecurityCenterContact(ctx, client, id.SecurityContactName, contact); err != nil {
			return fmt.Errorf("Creating Security Center Contact: %+v", err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id.SecurityContactName, contact); err != nil {
		return fmt.Errorf("Updating Security Center Contact: %+v", err)
	}

	return resourceSecurityCenterContactRead(d, meta)
}

func resourceSecurityCenterContactRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContactID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.SecurityContactName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if props := resp.ContactProperties; props != nil {
		d.Set("email", props.Email)
		d.Set("phone", props.Phone)
		d.Set("alert_notifications", props.AlertNotifications == security.On)
		d.Set("alerts_to_admins", props.AlertsToAdmins == security.AlertsToAdminsOn)
	}

	return nil
}

func resourceSecurityCenterContactDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContactID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.SecurityContactName); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
