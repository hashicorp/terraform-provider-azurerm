package privatedns

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsAaaaRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsAaaaRecordCreateUpdate,
		Read:   resourceArmPrivateDnsAaaaRecordRead,
		Update: resourceArmPrivateDnsAaaaRecordCreateUpdate,
		Delete: resourceArmPrivateDnsAaaaRecordDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AaaaRecordID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"records": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

func resourceArmPrivateDnsAaaaRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, privatedns.AAAA, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private DNS AAAA Record %q (Private Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_aaaa_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.RecordSet{
		Name: &name,
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata:    tags.Expand(t),
			TTL:         &ttl,
			AaaaRecords: expandAzureRmPrivateDnsAaaaRecords(d),
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, privatedns.AAAA, name, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating Private DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, privatedns.AAAA, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Private DNS AAAA Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Private DNS AAAA Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDnsAaaaRecordRead(d, meta)
}

func resourceArmPrivateDnsAaaaRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, privatedns.AAAA, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private DNS AAAA record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if err := d.Set("records", flattenAzureRmPrivateDnsAaaaRecords(resp.AaaaRecords)); err != nil {
		return err
	}
	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmPrivateDnsAaaaRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["AAAA"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Delete(ctx, resGroup, zoneName, privatedns.AAAA, name, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting Private DNS AAAA Record %s: %+v", name, err)
	}

	return nil
}

func flattenAzureRmPrivateDnsAaaaRecords(records *[]privatedns.AaaaRecord) []string {
	results := make([]string, 0)
	if records == nil {
		return results
	}

	for _, record := range *records {
		if record.Ipv6Address == nil {
			continue
		}

		results = append(results, *record.Ipv6Address)
	}

	return results
}

func expandAzureRmPrivateDnsAaaaRecords(d *schema.ResourceData) *[]privatedns.AaaaRecord {
	recordStrings := d.Get("records").(*schema.Set).List()
	records := make([]privatedns.AaaaRecord, len(recordStrings))

	for i, v := range recordStrings {
		ipv6 := v.(string)
		records[i] = privatedns.AaaaRecord{
			Ipv6Address: &ipv6,
		}
	}

	return &records
}
