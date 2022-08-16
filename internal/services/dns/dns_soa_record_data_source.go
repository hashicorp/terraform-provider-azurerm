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

			"tags": tags.Schema(),
		},
	}
}

func dataSourceDnsSoaRecordRead(d *pluginsdk.ResourceData, meta interface{}) error {
	recordSetsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	resourceGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	resp, err := recordSetsClient.Get(ctx, resourceGroup, zoneName, "@", dns.SOA)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: DNS SOA record (zone %s) was not found", zoneName)
		}
		return fmt.Errorf("reading DNS SOA record (zone %s): %+v", zoneName, err)
	}

	resourceId := parse.NewSoaRecordID(subscriptionId, resourceGroup, zoneName, "@")
	d.SetId(resourceId.ID())

	d.Set("resource_group_name", resourceGroup)
	d.Set("zone_name", zoneName)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if soaRecord := resp.SoaRecord; soaRecord != nil {
		d.Set("email", soaRecord.Email)
		d.Set("host_name", soaRecord.Host)
		d.Set("expire_time", soaRecord.ExpireTime)
		d.Set("minimum_ttl", soaRecord.MinimumTTL)
		d.Set("refresh_time", soaRecord.RefreshTime)
		d.Set("retry_time", soaRecord.RetryTime)
		d.Set("serial_number", soaRecord.SerialNumber)
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}
