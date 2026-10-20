// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	eventhubsClient := meta.(*clients.Client).Eventhub.EventHubsClient
	authorizationRulesClient := meta.(*clients.Client).Eventhub.EventHubAuthorizationRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := eventhubs.NewEventhubAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("eventhub_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			return pluginsdk.NonRetryableError(fmt.Errorf("expected %s was not found", id))
		}

		d.SetId(id.ID())

		if err := resourceEventHubAuthorizationRuleRead(d, meta); err != nil {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	})
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
