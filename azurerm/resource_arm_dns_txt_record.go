package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"record": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsTxtRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})

	records, err := expandAzureRmDnsTxtRecords(d)
	if err != nil {
		return err
	}

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:   expandTags(tags),
			TTL:        &ttl,
			TxtRecords: &records,
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.TXT, parameters, eTag, ifNoneMatch)
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
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["TXT"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.TXT)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS TXT record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("record", flattenAzureRmDnsTxtRecords(resp.TxtRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsTxtRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["TXT"]
	zoneName := id.Path["dnszones"]

	resp, error := client.Delete(ctx, resGroup, zoneName, name, dns.TXT, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS TXT Record %s: %+v", name, error)
	}

	return nil
}

func flattenAzureRmDnsTxtRecords(records *[]dns.TxtRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			txtRecord := make(map[string]interface{})

			if v := record.Value; v != nil {
				value := (*v)[0]
				txtRecord["value"] = value
			}

			results = append(results, txtRecord)
		}
	}

	return results
}

func expandAzureRmDnsTxtRecords(d *schema.ResourceData) ([]dns.TxtRecord, error) {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.TxtRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		value := []string{record["value"].(string)}

		txtRecord := dns.TxtRecord{
			Value: &value,
		}

		records[i] = txtRecord
	}

	return records, nil
}
