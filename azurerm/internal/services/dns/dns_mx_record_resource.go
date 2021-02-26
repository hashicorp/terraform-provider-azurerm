package dns

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDnsMxRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsMxRecordCreateUpdate,
		Read:   resourceDnsMxRecordRead,
		Update: resourceDnsMxRecordCreateUpdate,
		Delete: resourceDnsMxRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MxRecordID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "@",
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
				Set: resourceDnsMxRecordHash,
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

func resourceDnsMxRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.MX)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS MX Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_mx_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  tags.Expand(t),
			TTL:       &ttl,
			MxRecords: expandAzureRmDnsMxRecords(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.MX, parameters, "", ""); err != nil {
		return fmt.Errorf("Error creating/updating DNS MX Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.MX)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS MX Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS MX Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceDnsMxRecordRead(d, meta)
}

func resourceDnsMxRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MxRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DnszoneName, id.MXName, dns.MX)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS MX record %s: %v", id.MXName, err)
	}

	d.Set("name", id.MXName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.DnszoneName)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmDnsMxRecords(resp.MxRecords)); err != nil {
		return err
	}
	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceDnsMxRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MxRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.DnszoneName, id.MXName, dns.MX, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS MX Record %s: %+v", id.MXName, err)
	}

	return nil
}

// flatten creates an array of map where preference is a string to suit
// the expectations of the ResourceData schema, so that this data can be
// managed by Terradata state.
func flattenAzureRmDnsMxRecords(records *[]dns.MxRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

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
func expandAzureRmDnsMxRecords(d *schema.ResourceData) *[]dns.MxRecord {
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

	return &records
}

func resourceDnsMxRecordHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["preference"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["exchange"].(string)))
	}

	return schema.HashString(buf.String())
}
