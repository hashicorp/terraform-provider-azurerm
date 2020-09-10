package privatedns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsCNameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsCNameRecordCreateUpdate,
		Read:   resourceArmPrivateDnsCNameRecordRead,
		Update: resourceArmPrivateDnsCNameRecordCreateUpdate,
		Delete: resourceArmPrivateDnsCNameRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				// lower-cased due to the broken API https://github.com/Azure/azure-rest-api-specs/issues/6641
				ValidateFunc: validate.LowerCasedString,
			},

			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/6641
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"zone_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"record": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateDnsCNameRecordCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	zoneName := d.Get("zone_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, zoneName, privatedns.CNAME, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private DNS CNAME Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_cname_record", *existing.ID)
		}
	}

	ttl := int64(d.Get("ttl").(int))
	record := d.Get("record").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.RecordSet{
		Name: &name,
		RecordSetProperties: &privatedns.RecordSetProperties{
			Metadata: tags.Expand(t),
			TTL:      &ttl,
			CnameRecord: &privatedns.CnameRecord{
				Cname: &record,
			},
		},
	}

	eTag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	if _, err := client.CreateOrUpdate(ctx, resGroup, zoneName, privatedns.CNAME, name, parameters, eTag, ifNoneMatch); err != nil {
		return fmt.Errorf("Error creating/updating Private DNS CNAME Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, zoneName, privatedns.CNAME, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Private DNS CNAME Record %q (Zone %q / Resource Group %q): %s", name, zoneName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Private DNS CNAME Record %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDnsCNameRecordRead(d, meta)
}

func resourceArmPrivateDnsCNameRecordRead(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CNAME"]
	zoneName := id.Path["privateDnsZones"]

	resp, err := dnsClient.Get(ctx, resGroup, zoneName, privatedns.CNAME, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private DNS CNAME record %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("zone_name", zoneName)
	d.Set("ttl", resp.TTL)
	d.Set("fqdn", resp.Fqdn)

	if props := resp.RecordSetProperties; props != nil {
		if record := props.CnameRecord; record != nil {
			d.Set("record", record.Cname)
		}
	}

	return tags.FlattenAndSet(d, resp.Metadata)
}

func resourceArmPrivateDnsCNameRecordDelete(d *schema.ResourceData, meta interface{}) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["CNAME"]
	zoneName := id.Path["privateDnsZones"]

	_, err = dnsClient.Delete(ctx, resGroup, zoneName, privatedns.CNAME, name, "")
	if err != nil {
		return fmt.Errorf("Error deleting Private DNS CNAME Record %s: %+v", name, err)
	}

	return nil
}
