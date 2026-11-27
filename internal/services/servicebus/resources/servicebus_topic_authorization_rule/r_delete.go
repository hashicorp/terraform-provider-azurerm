// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusTopicAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).TopicsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := topics.ParseTopicAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.DeleteAuthorizationRule(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	if err := helpers.WaitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for Service Bus Namespace Disaster Recovery Configs (Namespace %q / Resource Group %q): %s", id.NamespaceName, id.ResourceGroupName, err)
	}

	return nil
}
