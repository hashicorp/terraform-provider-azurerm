package azurerm

import (
	"fmt"

	"time"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateDnsZoneCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.PrivateZonesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
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
	tags := d.Get("tags").(map[string]interface{})

	parameters := privatedns.PrivateZone{
		Location: &location,
		Tags:     expandTags(tags),
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation

	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("error creating/updating Private DNS Zone %q (Resource Group %q): %s", name, resGroup, err)
	}

	time.Sleep(time.Second * 10) // resource is slow to create / update, retry covers the create, sleeping to make update more reliable
	if err := resource.Retry(120*time.Second, retryPrivateDnsZonesClientGet(resGroup, name, meta)); err != nil {
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

func retryPrivateDnsZonesClientGet(resGroup string, name string, meta interface{}) func() *resource.RetryError {
	return func() *resource.RetryError {
		client := meta.(*ArmClient).privateDns.PrivateZonesClient
		ctx := meta.(*ArmClient).StopContext

		if _, err := client.Get(ctx, resGroup, name); err != nil {
			return resource.RetryableError(err)
		}

		return nil
	}
}

func resourceArmPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.PrivateZonesClient
	ctx := meta.(*ArmClient).StopContext

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
	d.Set("number_of_record_sets", resp.NumberOfRecordSets)
	d.Set("max_number_of_record_sets", resp.MaxNumberOfRecordSets)
	d.Set("max_number_of_virtual_network_links", resp.MaxNumberOfVirtualNetworkLinks)
	d.Set("max_number_of_virtual_network_links_with_registration", resp.MaxNumberOfVirtualNetworkLinksWithRegistration)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).privateDns.PrivateZonesClient
	ctx := meta.(*ArmClient).StopContext

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
