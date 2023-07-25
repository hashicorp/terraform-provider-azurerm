// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusTopicAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceServiceBusTopicAuthorizationRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AuthorizationRuleName(),
			},

			"topic_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: topics.ValidateTopicID,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.NamespaceName,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"queue_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.QueueName(),
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: resourcegroups.ValidateName,
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"topic_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.TopicName(),
				AtLeastOneOf: []string{"topic_id", "resource_group_name", "namespace_name", "topic_name"},
			},

			"listen": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"send": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"manage": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceServiceBusTopicAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsAuthClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var rgName string
	var nsName string
	var topicName string
	if v, ok := d.Get("topic_id").(string); ok && v != "" {
		topicId, err := topicsauthorizationrule.ParseTopicID(v)
		if err != nil {
			return fmt.Errorf("parsing topic ID %q: %+v", v, err)
		}
		rgName = topicId.ResourceGroupName
		nsName = topicId.NamespaceName
		topicName = topicId.TopicName
	} else {
		rgName = d.Get("resource_group_name").(string)
		nsName = d.Get("namespace_name").(string)
		topicName = d.Get("topic_name").(string)
	}

	id := topicsauthorizationrule.NewTopicAuthorizationRuleID(subscriptionId, rgName, nsName, topicName, d.Get("name").(string))
	resp, err := client.TopicsGetAuthorizationRule(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	keysResp, err := client.TopicsListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AuthorizationRuleName)
	d.Set("topic_name", id.TopicName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenTopicAuthorizationRuleRights(&props.Rights)
			d.Set("listen", listen)
			d.Set("send", send)
			d.Set("manage", manage)
		}
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}
