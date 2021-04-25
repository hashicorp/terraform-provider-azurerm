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

func resourceDnsARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsARecordCreateUpdate,
		Read:   resourceDnsARecordRead,
		Update: resourceDnsARecordCreateUpdate,
		Delete: resourceDnsARecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ARecordID(id)
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

			"records": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
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
				ConflictsWith: []string{"records"},

				// TODO: switch ConflictsWith for ExactlyOneOf when the Provider SDK's updated
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceDnsARecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	resourceId := parse.NewARecordID(subscriptionId, resGroup, zoneName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.A)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS A Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_a_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})
	targetResourceId := d.Get("target_resource_id").(string)
	recordsRaw := d.Get("records").(*schema.Set).List()

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:       tags.Expand(t),
			TTL:            &ttl,
			ARecords:       expandAzureRmDnsARecords(recordsRaw),
			TargetResource: &dns.SubResource{},
		},
	}

	if targetResourceId != "" {
		parameters.RecordSetProperties.TargetResource.ID = utils.String(targetResourceId)
	}

	// TODO: this can be removed when the provider SDK is upgraded
	if targetResourceId == "" && len(recordsRaw) == 0 {
		return fmt.Errorf("One of either `records` or `target_resource_id` must be specified")
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.A, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS A Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceDnsARecordRead(d, meta)
}

func resourceDnsARecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ARecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AName, dns.A)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS A record %s: %+v", id.AName, err)
	}

	d.Set("name", id.AName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.DnszoneName)

	d.Set("fqdn", resp.Fqdn)
	d.Set("ttl", resp.TTL)

	if err := d.Set("records", flattenAzureRmDnsARecords(resp.ARecords)); err != nil {
		return fmt.Errorf("Error setting `records`: %+v", err)
	}

	targetResourceId := ""
	if resp.TargetResource != nil && resp.TargetResource.ID != nil {
		targetResourceId = *resp.TargetResource.ID
	}
	d.Set("target_resource_id", targetResourceId)

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceDnsARecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ARecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Delete(ctx, id.ResourceGroup, id.DnszoneName, id.AName, dns.A, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS A Record %s: %+v", id.AName, err)
	}

	return nil
}

func expandAzureRmDnsARecords(input []interface{}) *[]dns.ARecord {
	records := make([]dns.ARecord, len(input))

	for i, v := range input {
		ipv4 := v.(string)
		records[i] = dns.ARecord{
			Ipv4Address: &ipv4,
		}
	}

	return &records
}

func flattenAzureRmDnsARecords(records *[]dns.ARecord) []string {
	if records == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, record := range *records {
		if record.Ipv4Address == nil {
			continue
		}

		results = append(results, *record.Ipv4Address)
	}

	return results
}
