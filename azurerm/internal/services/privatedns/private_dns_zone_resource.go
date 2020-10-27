package privatedns

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsZoneCreateUpdate,
		Read:   resourceArmPrivateDnsZoneRead,
		Update: resourceArmPrivateDnsZoneCreateUpdate,
		Delete: resourceArmPrivateDnsZoneDelete,
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
			},

			"number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_number_of_virtual_network_links_with_registration": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_zone", *existing.ID)
		}
	}

	location := "global"
	t := d.Get("tags").(map[string]interface{})

	parameters := privatedns.PrivateZone{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("error creating/updating Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for Private DNS Zone %q to become available: %+v", name, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("error retrieving Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Private DNS Zone %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmPrivateDnsZoneRead(d, meta)
}

func resourceArmPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["privateDnsZones"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Private DNS Zone %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if props := resp.PrivateZoneProperties; props != nil {
		d.Set("number_of_record_sets", props.NumberOfRecordSets)
		d.Set("max_number_of_record_sets", props.MaxNumberOfRecordSets)
		d.Set("max_number_of_virtual_network_links", props.MaxNumberOfVirtualNetworkLinks)
		d.Set("max_number_of_virtual_network_links_with_registration", props.MaxNumberOfVirtualNetworkLinksWithRegistration)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PrivateDns.PrivateZonesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["privateDnsZones"]

	etag := ""
	future, err := client.Delete(ctx, resGroup, name, etag)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("error deleting Private DNS Zone %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}
