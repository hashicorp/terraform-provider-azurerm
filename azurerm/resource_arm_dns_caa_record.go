package azurerm

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsCaaRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsCaaRecordCreateUpdate,
		Read:   resourceArmDnsCaaRecordRead,
		Update: resourceArmDnsCaaRecordCreateUpdate,
		Delete: resourceArmDnsCaaRecordDelete,
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

			"record": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flags": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"tag": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"issue",
								"issuewild",
								"iodef",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmDnsCaaRecordHash,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsCaaRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.CAA)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS CAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_caa_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:   expandTags(tags),
			TTL:        &ttl,
			CaaRecords: expandAzureRmDnsCaaRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.CAA, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS CAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.CAA)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS CAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS CAA Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsCaaRecordRead(d, meta)
}

func resourceArmDnsCaaRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CAA"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.CAA)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS CAA record %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("record", flattenAzureRmDnsCaaRecords(resp.CaaRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsCaaRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dns.RecordSetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CAA"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Delete(ctx, resGroup, zoneName, name, dns.CAA, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS CAA Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmDnsCaaRecords(records *[]dns.CaaRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			results = append(results, map[string]interface{}{
				"flags": *record.Flags,
				"tag":   *record.Tag,
				"value": *record.Value,
			})
		}
	}

	return results
}

func expandAzureRmDnsCaaRecords(d *schema.ResourceData) *[]dns.CaaRecord {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.CaaRecord, len(recordStrings))

	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		flags := int32(record["flags"].(int))
		tag := record["tag"].(string)
		value := record["value"].(string)

		caaRecord := dns.CaaRecord{
			Flags: &flags,
			Tag:   &tag,
			Value: &value,
		}

		records[i] = caaRecord
	}

	return &records
}

func resourceArmDnsCaaRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["flags"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["tag"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["value"].(string)))
	}

	return hashcode.String(buf.String())
}
