// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePrivateDnsMxRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateDnsMxRecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zone_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"record": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"preference": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"exchange": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourcePrivateDnsMxRecordHash,
			},

			"ttl": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePrivateDnsMxRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeMX, d.Get("name").(string))

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
	d.Set("zone_name", id.PrivateDnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if err := d.Set("record", flattenAzureRmPrivateDnsMxRecords(props.MxRecords)); err != nil {
				return err
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}
