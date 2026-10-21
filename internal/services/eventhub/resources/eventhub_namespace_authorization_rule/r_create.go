// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespaceAuthorizationRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := authorizationrulesnamespaces.NewAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.NamespacesGetAuthorizationRule(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_eventhub_namespace_authorization_rule", id.ID())
			}
		}
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	parameters := authorizationrulesnamespaces.AuthorizationRule{
		Name: &id.AuthorizationRuleName,
		Properties: &authorizationrulesnamespaces.AuthorizationRuleProperties{
			Rights: expandEventHubNamespaceAuthorizationRuleRights(d),
		},
	}

	if _, err := client.NamespacesCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}
	return resourceEventHubNamespaceAuthorizationRuleRead(d, meta)
}

func expandEventHubNamespaceAuthorizationRuleRights(d *pluginsdk.ResourceData) []authorizationrulesnamespaces.AccessRights {
	rights := make([]authorizationrulesnamespaces.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, authorizationrulesnamespaces.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, authorizationrulesnamespaces.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, authorizationrulesnamespaces.AccessRightsManage)
	}

	return rights
}
