// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusQueueAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceServiceBusQueueAuthorizationRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AuthorizationRuleName(),
			},

			"queue_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: queues.ValidateQueueID,
				AtLeastOneOf: []string{"queue_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.NamespaceName,
				AtLeastOneOf: []string{"queue_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"queue_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.QueueName(),
				AtLeastOneOf: []string{"queue_id", "resource_group_name", "namespace_name", "queue_name"},
			},

			"resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: resourcegroups.ValidateName,
				AtLeastOneOf: []string{"queue_id", "resource_group_name", "namespace_name", "queue_name"},
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

func dataSourceServiceBusQueueAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesAuthClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var rgName string
	var nsName string
	var queueName string
	if v, ok := d.Get("queue_id").(string); ok && v != "" {
		queueId, err := queuesauthorizationrule.ParseQueueID(v)
		if err != nil {
			return fmt.Errorf("parsing topic ID %q: %+v", v, err)
		}
		rgName = queueId.ResourceGroupName
		nsName = queueId.NamespaceName
		queueName = queueId.QueueName
	} else {
		rgName = d.Get("resource_group_name").(string)
		nsName = d.Get("namespace_name").(string)
		queueName = d.Get("queue_name").(string)
	}

	id := queuesauthorizationrule.NewQueueAuthorizationRuleID(subscriptionId, rgName, nsName, queueName, d.Get("name").(string))
	resp, err := client.QueuesGetAuthorizationRule(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AuthorizationRuleName)
	d.Set("queue_name", id.QueueName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("queue_id", queuesauthorizationrule.NewQueueID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.QueueName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenQueueAuthorizationRuleRights(&props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.QueuesListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if keysModel := keysResp.Model; keysModel != nil {
		d.Set("primary_key", keysModel.PrimaryKey)
		d.Set("primary_connection_string", keysModel.PrimaryConnectionString)
		d.Set("secondary_key", keysModel.SecondaryKey)
		d.Set("secondary_connection_string", keysModel.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
	}

	return nil
}
