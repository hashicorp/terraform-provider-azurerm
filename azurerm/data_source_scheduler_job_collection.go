package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSchedulerJobCollection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSchedulerJobCollectionRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"tags": tagsForDataSourceSchema(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"quota": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						//max_job_occurrence doesn't seem to do anything and always remains empty

						"max_job_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"max_recurrence_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// API documentation states the MaxRecurrence.Interval "Gets or sets the interval between retries."
						// however it does appear it is the max interval allowed for recurrences
						"max_retry_interval": {
							Type:       schema.TypeInt,
							Deprecated: "Renamed to `max_recurrence_interval` to match azure",
							Computed:   true,
						},

						"max_recurrence_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmSchedulerJobCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	collection, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(collection.Response) {
			return fmt.Errorf("Error: Scheduler Job Collection %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*collection.ID)

	//standard properties
	d.Set("name", collection.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := collection.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	flattenAndSetTags(d, collection.Tags)

	//resource specific
	if properties := collection.Properties; properties != nil {
		if sku := properties.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}
		d.Set("state", string(properties.State))

		if err := d.Set("quota", flattenAzureArmSchedulerJobCollectionQuota(properties.Quota)); err != nil {
			return fmt.Errorf("Error flattening quota for Job Collection %q (Resource Group %q): %+v", collection.Name, resourceGroup, err)
		}
	}

	return nil
}
