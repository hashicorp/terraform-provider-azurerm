package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSearchService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSearchServiceCreateUpdate,
		Read:   resourceArmSearchServiceRead,
		Delete: resourceArmSearchServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(search.Free),
					string(search.Basic),
					string(search.Standard),
					string(search.Standard2),
					string(search.Standard3),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"partition_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tags": tagsForceNewSchema(),
		},
	}
}

func resourceArmSearchServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).searchServicesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroupName := d.Get("resource_group_name").(string)
	skuName := d.Get("sku").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties := search.Service{
		Location: utils.String(location),
		Sku: &search.Sku{
			Name: search.SkuName(skuName),
		},
		ServiceProperties: &search.ServiceProperties{},
		Tags:              expandTags(tags),
	}

	if v, ok := d.GetOk("replica_count"); ok {
		replicaCount := int32(v.(int))
		properties.ServiceProperties.ReplicaCount = utils.Int32(replicaCount)
	}

	if v, ok := d.GetOk("partition_count"); ok {
		partitionCount := int32(v.(int))
		properties.ServiceProperties.PartitionCount = utils.Int32(partitionCount)
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroupName, name, properties, nil)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroupName, name, nil)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmSearchServiceRead(d, meta)
}

func resourceArmSearchServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).searchServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["searchServices"]

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Search Service %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Search Service: %+v", err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.ServiceProperties; props != nil {
		if count := props.PartitionCount; count != nil {
			d.Set("partition_count", int(*count))
		}

		if count := props.ReplicaCount; count != nil {
			d.Set("replica_count", int(*count))
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSearchServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).searchServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["searchServices"]

	resp, err := client.Delete(ctx, resourceGroup, name, nil)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Search Service %q (resource group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
