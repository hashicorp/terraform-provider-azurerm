package dns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDnsMxRecord() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDnsMxRecordRead,

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
							// TODO: this should become an Int
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"exchange": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceDnsMxRecordHash,
			},

			"ttl": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceDnsMxRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	recordSetsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	resp, err := recordSetsClient.Get(ctx, resourceGroup, zoneName, name, dns.MX)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: DNS MX record %s: (zone %s) was not found", name, zoneName)
		}
		return fmt.Errorf("reading DNS MX record %s (zone %s): %+v", name, zoneName, err)
	}

	resourceId := parse.NewMxRecordID(subscriptionId, resourceGroup, zoneName, name)
	d.SetId(resourceId.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("zone_name", zoneName)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmDnsMxRecords(resp.MxRecords)); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}
