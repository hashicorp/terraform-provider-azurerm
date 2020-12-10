package dns

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDnsAAAARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsAaaaRecordCreateUpdate,
		Read:   resourceDnsAaaaRecordRead,
		Update: resourceDnsAaaaRecordCreateUpdate,
		Delete: resourceDnsAaaaRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AaaaRecordID(id)
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Set:           set.HashIPv6Address,
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

			"tags": tags.Schema(),

			"target_resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"records"},
			},
		},
	}
}

func resourceDnsAaaaRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.AAAA)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_aaaa_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})
	recordsRaw := d.Get("records").(*schema.Set).List()
	targetResourceId := d.Get("target_resource_id").(string)

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:       tags.Expand(t),
			TTL:            &ttl,
			AaaaRecords:    expandAzureRmDnsAaaaRecords(recordsRaw),
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
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.AAAA, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.AAAA)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error retrieving DNS AAAA Record %q (Zone %q / Resource Group %q): ID was nil", name, zoneName, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceDnsAaaaRecordRead(d, meta)
}

func resourceDnsAaaaRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AaaaRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AAAAName, dns.AAAA)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS AAAA record %s: %v", id.AAAAName, err)
	}

	d.Set("name", id.AAAAName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.DnszoneName)

	d.Set("fqdn", resp.Fqdn)
	d.Set("ttl", resp.TTL)

	if err := d.Set("records", flattenAzureRmDnsAaaaRecords(resp.AaaaRecords)); err != nil {
		return fmt.Errorf("Error setting `records`: %+v", err)
	}

	targetResourceId := ""
	if resp.TargetResource != nil && resp.TargetResource.ID != nil {
		targetResourceId = *resp.TargetResource.ID
	}
	d.Set("target_resource_id", targetResourceId)

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceDnsAaaaRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AaaaRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Delete(ctx, id.ResourceGroup, id.DnszoneName, id.AAAAName, dns.AAAA, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS AAAA Record %s: %+v", id.AAAAName, err)
	}

	return nil
}

func expandAzureRmDnsAaaaRecords(input []interface{}) *[]dns.AaaaRecord {
	records := make([]dns.AaaaRecord, len(input))

	for i, v := range input {
		ipv6 := utils.NormalizeIPv6Address(v)
		records[i] = dns.AaaaRecord{
			Ipv6Address: &ipv6,
		}
	}

	return &records
}

func flattenAzureRmDnsAaaaRecords(records *[]dns.AaaaRecord) []string {
	if records == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, record := range *records {
		if record.Ipv6Address == nil {
			continue
		}

		results = append(results, utils.NormalizeIPv6Address(*record.Ipv6Address))
	}
	return results
}
