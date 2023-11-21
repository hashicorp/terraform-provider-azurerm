// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubAuthorizationRuleCreateUpdate,
		Read:   resourceEventHubAuthorizationRuleRead,
		Update: resourceEventHubAuthorizationRuleCreateUpdate,
		Delete: resourceEventHubAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubs.ParseEventhubAuthorizationRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.EventHubAuthorizationRuleV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: eventHubAuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(eventHubAuthorizationRuleCustomizeDiff),
	}
}

func resourceEventHubAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	eventhubsClient := meta.(*clients.Client).Eventhub.EventHubsClient
	authorizationRulesClient := meta.(*clients.Client).Eventhub.EventHubAuthorizationRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM EventHub Authorization Rule creation.")

	id := eventhubs.NewEventhubAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("eventhub_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := eventhubsClient.GetAuthorizationRule(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_authorization_rule", id.ID())
		}
	}

	locks.ByName(id.EventhubName, eventHubResourceName)
	defer locks.UnlockByName(id.EventhubName, eventHubResourceName)

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	parameters := authorizationruleseventhubs.AuthorizationRule{
		Name: &id.AuthorizationRuleName,
		Properties: &authorizationruleseventhubs.AuthorizationRuleProperties{
			Rights: expandEventHubAuthorizationRuleRights(d),
		},
	}

	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		localId := authorizationruleseventhubs.NewEventhubAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName, id.AuthorizationRuleName)
		if _, err := authorizationRulesClient.EventHubsCreateOrUpdateAuthorizationRule(ctx, localId, parameters); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creating %s: %+v", id, err))
		}

		read, err := eventhubsClient.GetAuthorizationRule(ctx, id)
		if err != nil {
			if response.WasNotFound(read.HttpResponse) {
				return pluginsdk.RetryableError(fmt.Errorf("expected %s to be created but was in non existent state, retrying", id))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("Expected %s was not be found", id))
		}

		d.SetId(id.ID())

		if err := resourceEventHubAuthorizationRuleRead(d, meta); err != nil {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	})
}

func resourceEventHubAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	eventHubsClient := meta.(*clients.Client).Eventhub.EventHubsClient
	authorizationRulesClient := meta.(*clients.Client).Eventhub.EventHubAuthorizationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubs.ParseEventhubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := eventHubsClient.GetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("eventhub_name", id.EventhubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			listen, send, manage := flattenEventHubAuthorizationRuleRights(properties.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	localId := authorizationruleseventhubs.NewEventhubAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.EventhubName, id.AuthorizationRuleName)
	keysResp, err := authorizationRulesClient.EventHubsListKeys(ctx, localId)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}

func resourceEventHubAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	eventhubClient := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubs.ParseEventhubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.EventhubName, eventHubResourceName)
	defer locks.UnlockByName(id.EventhubName, eventHubResourceName)

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if resp, err := eventhubClient.DeleteAuthorizationRule(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandEventHubAuthorizationRuleRights(d *pluginsdk.ResourceData) []authorizationruleseventhubs.AccessRights {
	rights := make([]authorizationruleseventhubs.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, authorizationruleseventhubs.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, authorizationruleseventhubs.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, authorizationruleseventhubs.AccessRightsManage)
	}

	return rights
}

func flattenEventHubAuthorizationRuleRights(rights []eventhubs.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	for _, right := range rights {
		switch right {
		case eventhubs.AccessRightsListen:
			listen = true
		case eventhubs.AccessRightsSend:
			send = true
		case eventhubs.AccessRightsManage:
			manage = true
		default:
			log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
		}
	}

	return listen, send, manage
}
