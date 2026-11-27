// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic_authorization_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusTopicAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).TopicsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	var id topics.TopicAuthorizationRuleId
	if topicIdLit := d.Get("topic_id").(string); topicIdLit != "" {
		topicId, err := topics.ParseTopicID(topicIdLit)
		if err != nil {
			return err
		}
		id = topics.NewTopicAuthorizationRuleID(topicId.SubscriptionId, topicId.ResourceGroupName, topicId.NamespaceName, topicId.TopicName, d.Get("name").(string))
	}

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAuthorizationRule(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_servicebus_topic_authorization_rule", id.ID())
			}
		}
	}

	parameters := topics.SBAuthorizationRule{
		Name: pointer.To(id.AuthorizationRuleName),
		Properties: &topics.SBAuthorizationRuleProperties{
			Rights: *expandTopicAuthorizationRuleRights(d),
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

	return resourceServiceBusTopicAuthorizationRuleRead(d, meta)
}

func expandTopicAuthorizationRuleRights(d *pluginsdk.ResourceData) *[]topics.AccessRights {
	rights := make([]topics.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, topics.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, topics.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, topics.AccessRightsManage)
	}

	return &rights
}
