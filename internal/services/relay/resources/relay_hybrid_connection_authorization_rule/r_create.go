// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	relayClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRelayHybridConnectionAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := relayClient.From(meta).HybridConnectionsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	resourceId := hybridconnections.NewHybridConnectionAuthorizationRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("hybrid_connection_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAuthorizationRule(ctx, resourceId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_relay_hybrid_connection_authorization_rule", resourceId.ID())
			}
		}
	}

	parameters := hybridconnections.AuthorizationRule{
		Name: pointer.To(resourceId.AuthorizationRuleName),
		Properties: &hybridconnections.AuthorizationRuleProperties{
			Rights: expandHybridConnectionAuthorizationRuleRights(d),
		},
	}

	if _, err := client.CreateOrUpdateAuthorizationRule(ctx, resourceId, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId.ID())
	}

	return resourceRelayHybridConnectionAuthorizationRuleRead(d, meta)
}
