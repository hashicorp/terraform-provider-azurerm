package dns

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDnsCNameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsCNameRecordCreateUpdate,
		Read:   resourceDnsCNameRecordRead,
		Update: resourceDnsCNameRecordCreateUpdate,
		Delete: resourceDnsCNameRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CnameRecordID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"record": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"target_resource_id"},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"record"},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDnsCNameRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	resourceId := parse.NewCnameRecordID(subscriptionId, resGroup, zoneName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.CNAME)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS CNAME Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_cname_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	record := d.Get("record").(string)
	t := d.Get("tags").(map[string]interface{})
	targetResourceId := d.Get("target_resource_id").(string)

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:       tags.Expand(t),
			TTL:            &ttl,
			CnameRecord:    &dns.CnameRecord{},
			TargetResource: &dns.SubResource{},
		},
	}

	if record != "" {
		parameters.RecordSetProperties.CnameRecord.Cname = utils.String(record)
	}

	if targetResourceId != "" {
		parameters.RecordSetProperties.TargetResource.ID = utils.String(targetResourceId)
	}

	// TODO: this can be removed when the provider SDK is upgraded
	if record == "" && targetResourceId == "" {
		return fmt.Errorf("One of either `record` or `target_resource_id` must be specified")
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.CNAME, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating CNAME Record %q (DNS Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceDnsCNameRecordRead(d, meta)
}

func resourceDnsCNameRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CnameRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.DnszoneName, id.CNAMEName, dns.CNAME)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving CNAME Record %s (DNS Zone %q / Resource Group %q): %+v", id.CNAMEName, id.DnszoneName, id.ResourceGroup, err)
	}

	d.Set("name", id.CNAMEName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.DnszoneName)

	d.Set("fqdn", resp.Fqdn)
	d.Set("ttl", resp.TTL)

	if props := resp.RecordSetProperties; props != nil {
		cname := ""
		if props.CnameRecord != nil && props.CnameRecord.Cname != nil {
			cname = *props.CnameRecord.Cname
		}
		d.Set("record", cname)

		targetResourceId := ""
		if props.TargetResource != nil && props.TargetResource.ID != nil {
			targetResourceId = *props.TargetResource.ID
		}
		d.Set("target_resource_id", targetResourceId)
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceDnsCNameRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CnameRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Delete(ctx, id.ResourceGroup, id.DnszoneName, id.CNAMEName, dns.CNAME, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting CNAME Record %q (DNS Zone %q / Resource Group %q): %+v", id.CNAMEName, id.DnszoneName, id.ResourceGroup, err)
	}

	return nil
}
