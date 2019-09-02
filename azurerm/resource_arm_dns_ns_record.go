package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsNsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsNsRecordCreateUpdate,
		Read:   resourceArmDnsNsRecordRead,
		Update: resourceArmDnsNsRecordCreateUpdate,
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDnsNsRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.NS)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS NS Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_ns_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  tags.Expand(t),
			TTL:       &ttl,
			NsRecords: expandAzureRmDnsNsRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.NS, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS NS Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.NS)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS NS Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS NS Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsNsRecordRead(d, meta)
}

func resourceArmDnsNsRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmDnsNsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["NS"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Delete(ctx, resGroup, zoneName, name, dns.NS, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS NS Record %s: %+v", name, err)
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

func expandAzureRmDnsNsRecords(d *schema.ResourceData) *[]dns.NsRecord {
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
	return &records
}
