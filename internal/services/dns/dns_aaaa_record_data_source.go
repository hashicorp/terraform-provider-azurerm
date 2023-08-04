// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDnsAAAARecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDnsAAAARecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zone_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"records": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      set.HashIPv6Address,
			},

			"ttl": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"target_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDnsAAAARecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	recordSetsClient := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeAAAA, d.Get("name").(string))

	resp, err := recordSetsClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.RelativeRecordSetName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.TTL)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("records", flattenAzureRmDnsAaaaRecords(props.AAAARecords)); err != nil {
				return fmt.Errorf("setting `records`: %+v", err)
			}

			targetResourceId := ""
			if props.TargetResource != nil && props.TargetResource.Id != nil {
				targetResourceId = *props.TargetResource.Id
			}
			d.Set("target_resource_id", targetResourceId)

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}
