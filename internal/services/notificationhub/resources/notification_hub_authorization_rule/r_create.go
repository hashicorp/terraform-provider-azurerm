// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hubs.NewNotificationHubAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("notification_hub_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.NotificationHubsGetAuthorizationRule(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_notification_hub_authorization_rule", id.ID())
			}
		}
	}

	locks.ByName(id.NotificationHubName, notificationHubResourceName)
	defer locks.UnlockByName(id.NotificationHubName, notificationHubResourceName)

	locks.ByName(id.NamespaceName, notificationHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, notificationHubNamespaceResourceName)

	manage := d.Get("manage").(bool)
	send := d.Get("send").(bool)
	listen := d.Get("listen").(bool)
	parameters := hubs.SharedAccessAuthorizationRuleResource{
		Properties: &hubs.SharedAccessAuthorizationRuleProperties{
			Rights: expandNotificationHubAuthorizationRuleRights(manage, send, listen),
		},
	}

	if _, err := client.NotificationHubsCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNotificationHubAuthorizationRuleRead(d, meta)
}

func expandNotificationHubAuthorizationRuleRights(manage bool, send bool, listen bool) []hubs.AccessRights {
	rights := make([]hubs.AccessRights, 0)

	if manage {
		rights = append(rights, hubs.AccessRightsManage)
	}

	if send {
		rights = append(rights, hubs.AccessRightsSend)
	}

	if listen {
		rights = append(rights, hubs.AccessRightsListen)
	}

	return rights
}
