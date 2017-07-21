package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsTxtRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsTxtRecordCreateOrUpdate,
		Read:   resourceArmDnsTxtRecordRead,
		Update: resourceArmDnsTxtRecordCreateOrUpdate,
		Delete: resourceArmDnsTxtRecordDelete,
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

			"record": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceArmDnsTxtRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	eTag := d.Get("etag").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	records, err := expandAzureRmDnsTxtRecords(d)
	props := dns.RecordSetProperties{
		Metadata:   metadata,
		TTL:        &ttl,
		TxtRecords: &records,
	}

	parameters := dns.RecordSet{
		Name:                &name,
		RecordSetProperties: &props,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
	resp, err := dnsClient.CreateOrUpdate(resGroup, zoneName, name, dns.TXT, parameters, eTag, "")
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS TXT Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsTxtRecordRead(d, meta)
}

func resourceArmDnsTxtRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["TXT"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(resGroup, zoneName, name, dns.TXT)
	if err != nil {
		return fmt.Errorf("Error reading DNS TXT record %s: %v", name, err)
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

	if err := d.Set("record", flattenAzureRmDnsTxtRecords(resp.TxtRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsTxtRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["TXT"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(resGroup, zoneName, name, dns.TXT, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS TXT Record %s: %s", name, error)
	}

	return nil
}

func flattenAzureRmDnsTxtRecords(records *[]dns.TxtRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			txtRecord := make(map[string]interface{})

			value := (*record.Value)[0]
			txtRecord["value"] = value

			results = append(results, txtRecord)
		}
	}

	return results
}

func expandAzureRmDnsTxtRecords(d *schema.ResourceData) ([]dns.TxtRecord, error) {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.TxtRecord, len(recordStrings))

	for i, v := range recordStrings {
		value := make([]string, 1)
		record := v.(map[string]interface{})
		value[0] = record["value"].(string)

		txtRecord := dns.TxtRecord{
			Value: &value,
		}

		records[i] = txtRecord
	}

	return records, nil
}
