// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_soa_record

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Private DNS SOA Record within Azure DNS.",
		Read:        dataSourcePrivateDnsSoaRecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Private DNS SOA Record.",
				Type:        pluginsdk.TypeString,
				Optional:    true,
				Default:     "@",
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zone_name": {
				Description: "The Private DNS Zone where the resource exists.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"email": {
				Description: "The email contact for the SOA record.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"host_name": {
				Description: "The domain name of the authoritative name server for the SOA record.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"expire_time": {
				Description: "The expire time for the SOA record.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"minimum_ttl": {
				Description: "The minimum Time To Live for the SOA record. By convention, it is used to determine the negative caching duration.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"refresh_time": {
				Description: "The refresh time for the SOA record.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"retry_time": {
				Description: "The retry time for the SOA record.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"serial_number": {
				Description: "The serial number for the SOA record.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"ttl": {
				Description: "The Time To Live (TTL) of the Private DNS record in seconds.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"fqdn": {
				Description: "The FQDN of the Private DNS SOA Record.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePrivateDnsSoaRecordRead(d *pluginsdk.ResourceData, meta any) error {
	recordSetsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := privatedns.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), privatedns.RecordTypeSOA, "@")
	resp, err := recordSetsClient.RecordSetsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("zone_name", id.PrivateDnsZoneName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ttl", props.Ttl)
			d.Set("fqdn", props.Fqdn)

			if soaRecord := props.SoaRecord; soaRecord != nil {
				d.Set("email", soaRecord.Email)
				d.Set("host_name", soaRecord.Host)
				d.Set("expire_time", soaRecord.ExpireTime)
				d.Set("minimum_ttl", soaRecord.MinimumTtl)
				d.Set("refresh_time", soaRecord.RefreshTime)
				d.Set("retry_time", soaRecord.RetryTime)
				d.Set("serial_number", soaRecord.SerialNumber)
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}
