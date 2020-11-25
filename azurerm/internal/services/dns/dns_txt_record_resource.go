package dns

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsTxtRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsTxtRecordCreateUpdate,
		Read:   resourceArmDnsTxtRecordRead,
		Update: resourceArmDnsTxtRecordCreateUpdate,
		Delete: resourceArmDnsTxtRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TxtRecordID(id)
			return err
		}),
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
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
					},
				},
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

func resourceArmDnsTxtRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, name, dns.TXT)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DNS TXT Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dns_txt_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:   tags.Expand(t),
			TTL:        &ttl,
			TxtRecords: expandAzureRmDnsTxtRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.TXT, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating DNS TXT Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, name, dns.TXT)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS TXT Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS TXT Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsTxtRecordRead(d, meta)
}

func resourceArmDnsTxtRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TxtRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DnszoneName, id.TXTName, dns.TXT)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS TXT record %s: %+v", id.TXTName, err)
	}

	d.Set("name", id.TXTName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.DnszoneName)
	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("record", flattenAzureRmDnsTxtRecords(resp.TxtRecords)); err != nil {
		return err
	}
	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmDnsTxtRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TxtRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.DnszoneName, id.TXTName, dns.TXT, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS TXT Record %s: %+v", id.TXTName, err)
	}

	return nil
}

func flattenAzureRmDnsTxtRecords(records *[]dns.TxtRecord) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	if records != nil {
		for _, record := range *records {
			txtRecord := make(map[string]interface{})

			if v := record.Value; v != nil {
				value := strings.Join(*v, "")
				txtRecord["value"] = value
			}

			results = append(results, txtRecord)
		}
	}

	return results
}

func expandAzureRmDnsTxtRecords(d *schema.ResourceData) *[]dns.TxtRecord {
	recordStrings := d.Get("record").(*schema.Set).List()
	records := make([]dns.TxtRecord, len(recordStrings))

	segmentLen := 254
	for i, v := range recordStrings {
		record := v.(map[string]interface{})
		v := record["value"].(string)

		var value []string
		for len(v) > segmentLen {
			value = append(value, v[:segmentLen])
			v = v[segmentLen:]
		}
		value = append(value, v)

		txtRecord := dns.TxtRecord{
			Value: &value,
		}

		records[i] = txtRecord
	}

	return &records
}
