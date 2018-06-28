package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDnsCNameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsCNameRecordCreateOrUpdate,
		Read:   resourceArmDnsCNameRecordRead,
		Update: resourceArmDnsCNameRecordCreateOrUpdate,
		Delete: resourceArmDnsCNameRecordDelete,
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
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Removed:  "Use `record` instead. This attribute will be removed in a future version",
			},

			"record": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsCNameRecordCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)
	ttl := int64(d.Get("ttl").(int))
	record := d.Get("record").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := dns.RecordSet{
		Name: &name,
		RecordSetProperties: &dns.RecordSetProperties{
			Metadata: expandTags(tags),
			TTL:      &ttl,
			CnameRecord: &dns.CnameRecord{
				Cname: &record,
			},
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := dnsClient.CreateOrUpdate(ctx, resGroup, zoneName, name, dns.CNAME, parameters, eTag, ifNoneMatch)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS CNAME Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsCNameRecordRead(d, meta)
}

func resourceArmDnsCNameRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CNAME"]
	zoneName := id.Path["dnszones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, name, dns.CNAME)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS CNAME record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)

	if props := resp.RecordSetProperties; props != nil {
		if record := props.CnameRecord; record != nil {
			d.Set("record", record.Cname)
		}
	}

	flattenAndSetTags(d, resp.Metadata)

	return nil
}

func resourceArmDnsCNameRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*ArmClient).dnsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CNAME"]
	zoneName := id.Path["dnszones"]

	resp, error := dnsClient.Delete(ctx, resGroup, zoneName, name, dns.CNAME, "")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting DNS CNAME Record %s: %+v", name, error)
	}

	return nil
}
