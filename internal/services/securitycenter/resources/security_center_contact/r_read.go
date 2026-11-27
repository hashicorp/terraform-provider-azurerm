// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_contact

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterContactRead(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).ContactsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.ContactID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.SecurityContactName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.SecurityContactName)

	if props := resp.ContactProperties; props != nil {
		d.Set("email", props.Email)
		d.Set("phone", props.Phone)
		d.Set("alert_notifications", props.AlertNotifications == security.On)
		d.Set("alerts_to_admins", props.AlertsToAdmins == security.AlertsToAdminsOn)
	}

	return nil
}
