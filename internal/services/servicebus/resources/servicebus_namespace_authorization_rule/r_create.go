// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).NamespacesAuthClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	namespaceAuthId, err := namespacesauthorizationrule.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	id := namespacesauthorizationrule.NewAuthorizationRuleID(namespaceAuthId.SubscriptionId, namespaceAuthId.ResourceGroupName, namespaceAuthId.NamespaceName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.NamespacesGetAuthorizationRule(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_servicebus_namespace_authorization_rule", id.ID())
			}
		}
	}

	parameters := namespacesauthorizationrule.SBAuthorizationRule{
		Name: pointer.To(id.AuthorizationRuleName),
		Properties: &namespacesauthorizationrule.SBAuthorizationRuleProperties{
			Rights: *expandAuthorizationRuleRights(d),
		},
	}

	if _, err := client.NamespacesCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := helpers.WaitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusNamespaceAuthorizationRuleRead(d, meta)
}
