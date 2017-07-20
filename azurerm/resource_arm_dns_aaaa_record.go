package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsAAAARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsAaaaRecordCreateOrUpdate,
		Read:   resourceArmDnsAaaaRecordRead,
		Update: resourceArmDnsAaaaRecordCreateOrUpdate,
		Delete: resourceArmDnsAaaaRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"records": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsAaaaRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	eTag := d.Get("etag").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	records, err := expandAzureRmDnsAaaaRecords(d)
	props := dns.RecordSetProperties{
		Metadata:    metadata,
		TTL:         &ttl,
		AaaaRecords: &records,
	}

	parameters := dns.RecordSet{
		Name:                &name,
		RecordSetProperties: &props,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
	resp, err := dnsClient.CreateOrUpdate(resGroup, zoneName, name, dns.AAAA, parameters, eTag, "")
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS AAAA Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsAaaaRecordRead(d, meta)
}

func resourceArmDnsAaaaRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(resGroup, zoneName, name, dns.AAAA)
	if err != nil {
		return fmt.Errorf("Error reading DNS AAAA record %s: %v", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)
	d.Set("etag", resp.Etag)

	if err := d.Set("records", flattenAzureRmDnsAaaaRecords(resp.AaaaRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsAaaaRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(resGroup, zoneName, name, dns.AAAA, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS AAAA Record %s: %s", name, error)
	}

	return nil
}

func flattenAzureRmDnsAaaaRecords(records *[]dns.AaaaRecord) []string {
	results := make([]string, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, *record.Ipv6Address)
		}
	}

	return results
}

func expandAzureRmDnsAaaaRecords(d *schema.ResourceData) ([]dns.AaaaRecord, error) {
	recordStrings := d.Get("records").(*schema.Set).List()
	records := make([]dns.AaaaRecord, len(recordStrings))

	for i, v := range recordStrings {
		ipv6 := v.(string)
		records[i] = dns.AaaaRecord{
			Ipv6Address: &ipv6,
		}
	}

	return records, nil
}
