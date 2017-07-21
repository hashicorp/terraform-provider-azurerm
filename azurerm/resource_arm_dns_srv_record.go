package azurerm

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsSrvRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsSrvRecordCreateOrUpdate,
		Read:   resourceArmDnsSrvRecordRead,
		Update: resourceArmDnsSrvRecordCreateOrUpdate,
		Delete: resourceArmDnsSrvRecordDelete,
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
						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"weight": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"target": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmDnsSrvRecordHash,
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

func resourceArmDnsSrvRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	eTag := d.Get("etag").(string)

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	records, err := expandAzureRmDnsSrvRecords(d)
	props := dns.RecordSetProperties{
		Metadata:   metadata,
		TTL:        &ttl,
		SrvRecords: &records,
	}

	parameters := dns.RecordSet{
		Name:                &name,
		RecordSetProperties: &props,
	}

	//last parameter is set to empty to allow updates to records after creation
	// (per SDK, set it to '*' to prevent updates, all other values are ignored)
	resp, err := dnsClient.CreateOrUpdate(resGroup, zoneName, name, dns.SRV, parameters, eTag, "")
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS SRV Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsSrvRecordRead(d, meta)
}

func resourceArmDnsSrvRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["SRV"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(resGroup, zoneName, name, dns.SRV)
	if err != nil {
		return fmt.Errorf("Error reading DNS SRV record %s: %v", name, err)
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

	if err := d.Set("record", flattenAzureRmDnsSrvRecords(resp.SrvRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsSrvRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["SRV"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(resGroup, zoneName, name, dns.SRV, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS SRV Record %s: %s", name, error)
	}

	return nil
}

func flattenAzureRmDnsSrvRecords(records *[]dns.SrvRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, map[string]interface{}{
				"priority": *record.Priority,
				"weight":   *record.Weight,
				"port":     *record.Port,
				"target":   *record.Target,
			})
		}
	}

	return results
}

func expandAzureRmDnsSrvRecords(d *schema.ResourceData) ([]dns.SrvRecord, error) {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.SrvRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		priority := int32(record["priority"].(int))
		weight := int32(record["weight"].(int))
		port := int32(record["port"].(int))
		target := record["target"].(string)

		srvRecord := dns.SrvRecord{
			Priority: &priority,
			Weight:   &weight,
			Port:     &port,
			Target:   &target,
		}

		records[i] = srvRecord
	}

	return records, nil
}

func resourceArmDnsSrvRecordHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["priority"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["weight"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["target"].(string)))

	return hashcode.String(buf.String())
}
