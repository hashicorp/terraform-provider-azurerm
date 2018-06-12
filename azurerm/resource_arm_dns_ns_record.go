package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsNsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsNsRecordCreateOrUpdate,
		Read:   resourceArmDnsNsRecordRead,
		Update: resourceArmDnsNsRecordCreateOrUpdate,
		Delete: resourceArmDnsNsRecordDelete,
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

			"records": {
				Type: schema.TypeList,
				//TODO: add `Required: true` once we remove the `record` attribute
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"record"},
			},

			"record": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Deprecated:    "This field has been replaced by `records`",
				ConflictsWith: []string{"records"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nsdname": {
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

func resourceArmDnsNsRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})
	records, err := expandAzureRmDnsNsRecords(d)
	if err != nil {
		return err
	}

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  expandTags(tags),
			TTL:       &ttl,
			NsRecords: &records,
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := dnsClient.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.NS, parameters, eTag, ifNoneMatch)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS NS Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsNsRecordRead(d, meta)
}

func resourceArmDnsNsRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["NS"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, name, dns.NS)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS NS record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("records", flattenAzureRmDnsNsRecords(resp.NsRecords)); err != nil {
		return fmt.Errorf("Error settings `records`: %+v", err)
	}

	//TODO: remove this once we remove the `record` attribute
	if err := d.Set("record", flattenAzureRmDnsNsRecordsSet(resp.NsRecords)); err != nil {
		return fmt.Errorf("Error settings `record`: %+v", err)
	}

	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsNsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["NS"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(ctx, resGroup, zoneName, name, dns.NS, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS NS Record %s: %+v", name, error)
	}

	return nil
}

//TODO: remove this once we remove the `record` attribute
func flattenAzureRmDnsNsRecordsSet(records *[]dns.NsRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			nsRecord := make(map[string]interface{})
			nsRecord["nsdname"] = *record.Nsdname
			results = append(results, nsRecord)
		}
	}

	return results
}

func flattenAzureRmDnsNsRecords(records *[]dns.NsRecord) []string {
	results := make([]string, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, *record.Nsdname)
		}
	}

	return results
}

func expandAzureRmDnsNsRecords(d *schema.ResourceData) ([]dns.NsRecord, error) {
	var records []dns.NsRecord

	//TODO: remove this once we remove the `record` attribute
	if d.HasChange("records") || !d.HasChange("record") {
		recordStrings := d.Get("records").([]interface{})
		records = make([]dns.NsRecord, len(recordStrings))
		for i, v := range recordStrings {
			record := v.(string)

			nsRecord := dns.NsRecord{
				Nsdname: &record,
			}

			records[i] = nsRecord
		}
	} else {
		recordList := d.Get("record").(*schema.Set).List()
		if len(recordList) != 0 {
			records = make([]dns.NsRecord, len(recordList))
			for i, v := range recordList {
				record := v.(map[string]interface{})
				nsdname := record["nsdname"].(string)
				nsRecord := dns.NsRecord{
					Nsdname: &nsdname,
				}

				records[i] = nsRecord
			}
		}
	}
	return records, nil
}
