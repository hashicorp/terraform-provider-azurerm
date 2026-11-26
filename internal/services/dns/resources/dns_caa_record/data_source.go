// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dns_caa_record

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

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing DNS CAA Record within Azure DNS.",
		Read:        dataSourceDnsCaaRecordRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the DNS CAA Record.",
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
				Description: "A list of values that make up the CAA record.",
				Type:        pluginsdk.TypeSet,
				Computed:    true,
				Set:         resourceDnsCaaRecordHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"flags": {
							Description: "Extensible CAA flags, currently only 1 is implemented to set the issuer critical flag.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"tag": {
							Description: "A property tag, options are `issue`, `issuewild`, `iodef`, and `contactemail`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"value": {
							Description: "A property value such as a registrar domain.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"ttl": {
				Description: "The Time To Live (TTL) of the DNS record in seconds.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"fqdn": {
				Description: "The FQDN of the DNS CAA Record.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDnsCaaRecordRead(d *pluginsdk.ResourceData, meta any) error {
	recordSetsClient := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := recordsets.NewRecordTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("zone_name").(string), recordsets.RecordTypeCAA, d.Get("name").(string))

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

			if err := d.Set("record", flattenAzureRmDnsCaaRecords(props.CaaRecords)); err != nil {
				return err
			}

			return tags.FlattenAndSet(d, props.Metadata)
		}
	}

	return nil
}
