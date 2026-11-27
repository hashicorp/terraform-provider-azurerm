// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusQueueAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.DeleteAuthorizationRule(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := helpers.WaitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", *id, err)
	}

	return nil
}
