package dns

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsNsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsNsRecordCreate,
		Read:   resourceArmDnsNsRecordRead,
		Update: resourceArmDnsNsRecordUpdate,
		Delete: resourceArmDnsNsRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DnsNsRecordID(id)
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
				ForceNew: true,
			},

			"records": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func resourceArmDnsNsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if features.ShouldResourcesBeImported() {
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

	recordsRaw := d.Get("records").([]interface{})
	records := expandAzureRmDnsNsRecords(recordsRaw)

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata:  tags.Expand(t),
			TTL:       &ttl,
			NsRecords: records,
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.NS, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating DNS NS Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
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

func resourceArmDnsNsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["NS"]
	zoneName := id.Path["dnszones"]

	existing, err := client.Get(ctx, resGroup, zoneName, name, dns.NS)
	if err != nil {
		return fmt.Errorf("Error retrieving NS %q (DNS Zone %q / Resource Group %q): %+v", name, zoneName, resGroup, err)
	}

	if existing.RecordSetProperties == nil {
		return fmt.Errorf("Error retrieving NS %q (DNS Zone %q / Resource Group %q): `properties` was nil", name, zoneName, resGroup)
	}

	if d.HasChange("records") {
		recordsRaw := d.Get("records").([]interface{})
		records := expandAzureRmDnsNsRecords(recordsRaw)
		existing.RecordSetProperties.NsRecords = records
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		existing.RecordSetProperties.Metadata = tags.Expand(t)
	}

	if d.HasChange("ttl") {
		existing.RecordSetProperties.TTL = utils.Int64(int64(d.Get("ttl").(int)))
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.NS, existing, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error updating DNS NS Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	return resourceArmDnsNsRecordRead(d, meta)
}

func resourceArmDnsNsRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsNsRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Get(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.NS)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS NS record %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zone_name", id.ZoneName)

	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if props := resp.RecordSetProperties; props != nil {
		if err := d.Set("records", flattenAzureRmDnsNsRecords(props.NsRecords)); err != nil {
			return fmt.Errorf("Error settings `records`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmDnsNsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).Dns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DnsNsRecordID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.Delete(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.NS, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS NS Record %s: %+v", id.Name, err)
	}

	return nil
}

func flattenAzureRmDnsNsRecords(records *[]dns.NsRecord) []interface{} {
	if records == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, record := range *records {
		if record.Nsdname == nil {
			continue
		}

		results = append(results, *record.Nsdname)
	}

	return results
}

func expandAzureRmDnsNsRecords(input []interface{}) *[]dns.NsRecord {
	records := make([]dns.NsRecord, len(input))
	for i, v := range input {
		record := v.(string)

		nsRecord := dns.NsRecord{
			Nsdname: &record,
		}

		records[i] = nsRecord
	}
	return &records
}
