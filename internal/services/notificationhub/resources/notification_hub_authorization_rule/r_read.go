// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_authorization_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubAuthorizationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hubs.ParseNotificationHubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NotificationHubsGetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := client.NotificationHubsListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing access keys for %s: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("notification_hub_name", id.NotificationHubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			manage, send, listen := flattenNotificationHubAuthorizationRuleRights(&props.Rights)
			d.Set("manage", manage)
			d.Set("send", send)
			d.Set("listen", listen)
		}
	}

	if keysModel := keysResp.Model; keysModel != nil {
		d.Set("primary_access_key", keysModel.PrimaryKey)
		d.Set("secondary_access_key", keysModel.SecondaryKey)
		d.Set("primary_connection_string", keysModel.PrimaryConnectionString)
		d.Set("secondary_connection_string", keysModel.SecondaryConnectionString)
	}

	return nil
}

func flattenNotificationHubAuthorizationRuleRights(input *[]hubs.AccessRights) (manage bool, send bool, listen bool) {
	if input == nil {
		return
	}

	for _, right := range *input {
		switch right {
		case hubs.AccessRightsManage:
			manage = true
			continue
		case hubs.AccessRightsSend:
			send = true
			continue
		case hubs.AccessRightsListen:
			listen = true
			continue
		}
	}

	return manage, send, listen
}
