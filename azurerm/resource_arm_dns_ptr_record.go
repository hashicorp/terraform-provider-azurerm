package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsPtrRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsPtrRecordCreate,
		Read:   resourceArmDnsPtrRecordRead,
		Update: resourceArmDnsPtrRecordCreate,
		Delete: resourceArmDnsPtrRecordDelete,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsPtrRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	dnsClient := client.dnsClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})

	metadata := expandTags(tags)

	recordStrings := d.Get("records").(*schema.Set).List()
	records := make([]dns.PtrRecord, len(recordStrings))

	for i, v := range recordStrings {
		fqdn := v.(string)
		records[i] = dns.PtrRecord{
			Ptrdname: &fqdn,
		}
	}

	props := dns.RecordSetProperties{
		Metadata:   metadata,
		TTL:        &ttl,
		PtrRecords: &records,
	}

	parameters := dns.RecordSet{
		Name:                &name,
		RecordSetProperties: &props,
	}

	_, err := dnsClient.CreateOrUpdate(resGroup, zoneName, name, dns.PTR, parameters, "", "")
	if err != nil {
		return err
	}

	rec, err := dnsClient.Get(resGroup, zoneName, name, dns.PTR)
	if err != nil {
		return err
	}

	if rec.ID == nil {
		return fmt.Errorf("Cannot read DNS PTR Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*rec.ID)

	return resourceArmDnsPtrRecordRead(d, meta)
}

func resourceArmDnsPtrRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	dnsClient := client.dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["PTR"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(resGroup, zoneName, name, dns.PTR)
	if err != nil {
		return fmt.Errorf("Error reading DNS PTR record %s: %s", name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zonename", zoneName)
	d.Set("ttl", resp.TTL)

	if resp.PtrRecords != nil {
		records := make([]string, 0, len(*resp.PtrRecords))
		for _, record := range *resp.PtrRecords {
			records = append(records, *record.Ptrdname)
		}

		if err := d.Set("records", records); err != nil {
			return err
		}
	}

	flattenAndSetTags(d, resp.Metadata)
	return nil
}

func resourceArmDnsPtrRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	dnsClient := client.dnsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["PTR"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(resGroup, zoneName, name, dns.PTR, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS PTR Record %s: %s", name, error)
	}

	return nil
}
