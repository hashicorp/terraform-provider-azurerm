// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusQueueAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueId, err := queues.ParseQueueID(d.Get("queue_id").(string))
	if err != nil {
		return err
	}

	id := queues.NewQueueAuthorizationRuleID(queueId.SubscriptionId, queueId.ResourceGroupName, queueId.NamespaceName, queueId.QueueName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAuthorizationRule(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_servicebus_queue_authorization_rule", id.ID())
			}
		}
	}

	parameters := queues.SBAuthorizationRule{
		Name: pointer.To(id.AuthorizationRuleName),
		Properties: &queues.SBAuthorizationRuleProperties{
			Rights: *expandQueueAuthorizationRuleRights(d),
		},
	}

	if _, err := client.CreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := helpers.WaitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusQueueAuthorizationRuleRead(d, meta)
}

func expandQueueAuthorizationRuleRights(d *pluginsdk.ResourceData) *[]queues.AccessRights {
	rights := make([]queues.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, queues.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, queues.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, queues.AccessRightsManage)
	}

	return &rights
}
