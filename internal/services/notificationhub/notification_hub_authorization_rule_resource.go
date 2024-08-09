// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Read:   resourceNotificationHubAuthorizationRuleRead,
		Update: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Delete: resourceNotificationHubAuthorizationRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hubs.ParseNotificationHubAuthorizationRuleID(id)
			return err
		}),
		// TODO: customizeDiff for send+listen when manage selected

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NotificationHubAuthorizationRuleResourceV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"notification_hub_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"manage": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"listen": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: features.FourPointOhBeta(),
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: features.FourPointOhBeta(),
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceNotificationHubAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hubs.NewNotificationHubAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("notification_hub_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
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

func resourceNotificationHubAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceNotificationHubAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hubs.ParseNotificationHubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NotificationHubName, notificationHubResourceName)
	defer locks.UnlockByName(id.NotificationHubName, notificationHubResourceName)

	locks.ByName(id.NamespaceName, notificationHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, notificationHubNamespaceResourceName)

	resp, err := client.NotificationHubsDeleteAuthorizationRule(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
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
