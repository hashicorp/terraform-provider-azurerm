// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRelayNamespaceAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Relay.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := namespaces.NewAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAuthorizationRule(ctx, resourceId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_relay_namespace_authorization_rule", resourceId.ID())
			}
		}
	}

	parameters := namespaces.AuthorizationRule{
		Name: pointer.To(resourceId.AuthorizationRuleName),
		Properties: &namespaces.AuthorizationRuleProperties{
			Rights: expandAuthorizationRuleRights(d),
		},
	}

	if _, err := client.CreateOrUpdateAuthorizationRule(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId.ID())
	}

	return resourceRelayNamespaceAuthorizationRuleRead(d, meta)
}
