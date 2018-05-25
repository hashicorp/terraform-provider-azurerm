package azurerm

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsMxRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsMxRecordCreateOrUpdate,
		Read:   resourceArmDnsMxRecordRead,
		Update: resourceArmDnsMxRecordCreateOrUpdate,
		Delete: resourceArmDnsMxRecordDelete,
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
						"preference": {
							// TODO: this should become an Int
							Type:     schema.TypeString,
							Required: true,
						},

						"exchange": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmDnsMxRecordHash,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsMxRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	tags := d.Get("tags").(map[string]interface{})
	records, err := expandAzureRmDnsMxRecords(d)
	if err != nil {
		return err
	}

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  expandTags(tags),
			TTL:       &ttl,
			MxRecords: &records,
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.MX, parameters, eTag, ifNoneMatch)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS MX Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsMxRecordRead(d, meta)
}

func resourceArmDnsMxRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["MX"]
	zoneName := id.Path["dnszones"]

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.MX)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS MX record %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if err := d.Set("record", flattenAzureRmDnsMxRecords(resp.MxRecords)); err != nil {
		return err
	}
	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsMxRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["MX"]
	zoneName := id.Path["dnszones"]

	resp, error := client.Delete(ctx, resGroup, zoneName, name, dns.MX, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS MX Record %s: %+v", name, error)
	}

	return nil
}

// flatten creates an array of map where preference is a string to suit
// the expectations of the ResourceData schema, so that this data can be
// managed by Terradata state.
func flattenAzureRmDnsMxRecords(records *[]dns.MxRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(*records))

	if records != nil {
		for _, record := range *records {
			preferenceI32 := *record.Preference
			preference := strconv.Itoa(int(preferenceI32))
			results = append(results, map[string]interface{}{
				"preference": preference,
				"exchange":   *record.Exchange,
			})
		}
	}

	return results
}

// expand creates an array of dns.MxRecord, that is, the array needed
// by azure-sdk-for-go to manipulate azure resources, hence Preference
// is an int32
func expandAzureRmDnsMxRecords(d *schema.ResourceData) ([]dns.MxRecord, error) {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.MxRecord, len(recordStrings))

	for i, v := range recordStrings {
		mxrecord := v.(map[string]interface{})
		preference := mxrecord["preference"].(string)
		i64, _ := strconv.ParseInt(preference, 10, 32)
		i32 := int32(i64)
		exchange := mxrecord["exchange"].(string)

		records[i] = dns.MxRecord{
			Preference: &i32,
			Exchange:   &exchange,
		}
	}

	return records, nil
}

func resourceArmDnsMxRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["preference"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["exchange"].(string)))
	}

	return hashcode.String(buf.String())
}
