package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSchedulerJobCollection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSchedulerJobCollectionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		DeprecationMessage: "Scheduler Job Collection has been deprecated in favour of Logic Apps - more information can be found at https://docs.microsoft.com/en-us/azure/scheduler/migrate-from-scheduler-to-logic-apps",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),

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
	client := meta.(*clients.Client).Scheduler.JobCollectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	collection, err := client.Get(ctx, resourceGroup, name) //nolint: megacheck
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	//resource specific
	if properties := collection.Properties; properties != nil {
		if sku := properties.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}
		d.Set("state", string(properties.State))

		if err := d.Set("quota", flattenAzureArmSchedulerJobCollectionQuota(properties.Quota)); err != nil {
			return fmt.Errorf("Error setting quota for Job Collection %q (Resource Group %q): %+v", *collection.Name, resourceGroup, err)
		}
	}

	return tags.FlattenAndSet(d, collection.Tags)
}
