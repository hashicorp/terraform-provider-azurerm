// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dns_cname_record

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing DNS CNAME Record within Azure DNS.",
		Read:        dataSourceDnsCNameRecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the DNS CNAME Record.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zone_name": {
				Description: "The DNS Zone where the resource exists.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"record": {
				Description: "The target of the CNAME.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ttl": {
				Description: "The Time To Live (TTL) of the DNS record in seconds.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"fqdn": {
				Description: "The FQDN of the DNS CNAME Record.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"target_resource_id": {
				Description: "The Azure resource ID of the target object from where the DNS resource value is taken.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDnsCNameRecordRead(d *pluginsdk.ResourceData, meta any) error {
	recordSetsClient := client.From(meta).RecordSets
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()
	subscriptionId := account.From(meta).SubscriptionId

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeCNAME, d.Get("name").(string))

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

			cname := ""
			if props.CNAMERecord != nil && props.CNAMERecord.Cname != nil {
				cname = *props.CNAMERecord.Cname
			}
			d.Set("record", cname)

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
