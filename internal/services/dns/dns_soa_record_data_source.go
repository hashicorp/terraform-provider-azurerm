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
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDnsSoaRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDnsSoaRecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "@",
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zone_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"email": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expire_time": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"minimum_ttl": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"refresh_time": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"retry_time": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"serial_number": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
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

func dataSourceDnsSoaRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	recordSetsClient := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeSOA, "@")
	resp, err := recordSetsClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.DnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.TTL)
			d.Set("fqdn", props.Fqdn)

			if soaRecord := props.SOARecord; soaRecord != nil {
				d.Set("email", soaRecord.Email)
				d.Set("host_name", soaRecord.Host)
				d.Set("expire_time", soaRecord.ExpireTime)
				d.Set("minimum_ttl", soaRecord.MinimumTTL)
				d.Set("refresh_time", soaRecord.RefreshTime)
				d.Set("retry_time", soaRecord.RetryTime)
				d.Set("serial_number", soaRecord.SerialNumber)
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}
