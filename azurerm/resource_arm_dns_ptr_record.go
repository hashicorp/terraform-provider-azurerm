package azurerm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsPtrRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsPtrRecordCreateUpdate,
		Read:   resourceArmDnsPtrRecordRead,
		Update: resourceArmDnsPtrRecordCreateUpdate,
		Delete: resourceArmDnsPtrRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
		},
	}
}

func resourceArmDnsPtrRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.PTR)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS PTR Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_ptr_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:   tags.Expand(t),
			TTL:        &ttl,
			PtrRecords: expandAzureRmDnsPtrRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.PTR, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS PTR Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.PTR)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS PTR Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS PTR Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsPtrRecordRead(d, meta)
}

func resourceArmDnsPtrRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	dnsClient := client.Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["PTR"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, name, dns.PTR)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading DNS PTR record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("records", flattenAzureRmDnsPtrRecords(resp.PtrRecords)); err != nil {
		return err
	}
	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmDnsPtrRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	dnsClient := client.Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["PTR"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Delete(ctx, resGroup, zoneName, name, dns.PTR, "")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error deleting DNS PTR Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmDnsPtrRecords(records *[]dns.PtrRecord) []string {
	results := make([]string, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, *record.Ptrdname)
		}
	}

	return results
}

func expandAzureRmDnsPtrRecords(d *schema.ResourceData) *[]dns.PtrRecord {
	recordStrings := d.Get("records").(*schema.Set).List()
	records := make([]dns.PtrRecord, len(recordStrings))

	for i, v := range recordStrings {
		fqdn := v.(string)
		records[i] = dns.PtrRecord{
			Ptrdname: &fqdn,
		}
	}

	return &records
}
