// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_contact

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterContactCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	// TODO: split this Create/Update
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewContactID(subscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.SecurityContactName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_security_center_contact", id.ID())
			}
		}
	}

	contact := security.Contact{
		ContactProperties: &security.ContactProperties{
			Email: pointer.To(d.Get("email").(string)),
			Phone: pointer.To(d.Get("phone").(string)),
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
			return fmt.Errorf("creating Security Center Contact: %+v", err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id.SecurityContactName, contact); err != nil {
		return fmt.Errorf("updating Security Center Contact: %+v", err)
	}

	return resourceSecurityCenterContactRead(d, meta)
}
