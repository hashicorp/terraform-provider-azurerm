package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/search"
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

	name := d.Get("name").(string)
	location := d.Get("location").(string)
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
		replica_count := int32(v.(int))
		properties.ServiceProperties.ReplicaCount = utils.Int32(replica_count)
	}

	if v, ok := d.GetOk("partition_count"); ok {
		partition_count := int32(v.(int))
		properties.ServiceProperties.PartitionCount = utils.Int32(partition_count)
	}

	_, err := client.CreateOrUpdate(resourceGroupName, name, properties, nil)
	if err != nil {
		return err
	}

	resp, err := client.Get(resourceGroupName, name, nil)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmSearchServiceRead(d, meta)
}

func resourceArmSearchServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).searchServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["searchServices"]

	resp, err := client.Get(resourceGroup, name, nil)
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
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if resp.Sku != nil {
		d.Set("sku", string(resp.Sku.Name))
	}

	if props := resp.ServiceProperties; props != nil {
		if props.PartitionCount != nil {
			d.Set("partition_count", int(*props.PartitionCount))
		}

		if props.ReplicaCount != nil {
			d.Set("replica_count", int(*props.ReplicaCount))
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSearchServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).searchServicesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["searchServices"]

	resp, err := client.Delete(resourceGroup, name, nil)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Search Service: %+v", err)
	}

	return nil
}
