// nolint: megacheck
// entire automation SDK has been depreciated in v21.3 in favor of logic apps, an entirely different service.
package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSchedulerJobCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSchedulerJobCollectionCreateUpdate,
		Read:   resourceArmSchedulerJobCollectionRead,
		Update: resourceArmSchedulerJobCollectionCreateUpdate,
		Delete: resourceArmSchedulerJobCollectionDelete,

		DeprecationMessage: "Scheduler Job Collection has been deprecated in favour of Logic Apps - more information can be found at https://docs.microsoft.com/en-us/azure/scheduler/migrate-from-scheduler-to-logic-apps",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-_a-zA-Z0-9]{0,99}$"),
					"Job Collection Name name must be 1 - 100 characters long, start with a letter and contain only letters, numbers, hyphens and underscores.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": tags.Schema(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(scheduler.Free),
					string(scheduler.Standard),
					string(scheduler.P10Premium),
					string(scheduler.P20Premium),
				}, true),
			},

			//optional
			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(scheduler.Enabled),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(scheduler.Enabled),
					string(scheduler.Suspended),
					string(scheduler.Disabled),
				}, true),
			},

			"quota": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						//max_job_occurrence doesn't seem to do anything and always remains empty

						"max_job_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"max_recurrence_frequency": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(scheduler.Minute),
								string(scheduler.Hour),
								string(scheduler.Day),
								string(scheduler.Week),
								string(scheduler.Month),
							}, true),
						},

						// API documentation states the MaxRecurrence.Interval "Gets or sets the interval between retries."
						// however it does appear it is the max interval allowed for recurrences
						"max_retry_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							Deprecated:   "Renamed to `max_recurrence_interval` to match azure",
							ValidateFunc: validation.IntAtLeast(1), //changes depending on the frequency, unknown maximums
						},

						"max_recurrence_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1), //changes depending on the frequency, unknown maximums
						},
					},
				},
			},
		},
	}
}

func resourceArmSchedulerJobCollectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).scheduler.JobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating/updating Scheduler Job Collection %q (resource group %q)", name, resourceGroup)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_scheduler_job_collection", *existing.ID)
		}
	}

	collection := scheduler.JobCollectionDefinition{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		Properties: &scheduler.JobCollectionProperties{
			Sku: &scheduler.Sku{
				Name: scheduler.SkuDefinition(d.Get("sku").(string)),
			},
		},
	}

	if state, ok := d.Get("state").(string); ok {
		collection.Properties.State = scheduler.JobCollectionState(state)
	}
	collection.Properties.Quota = expandAzureArmSchedulerJobCollectionQuota(d)

	//create job collection
	collection, err := client.CreateOrUpdate(ctx, resourceGroup, name, collection)
	if err != nil {
		return fmt.Errorf("Error creating/updating Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//ensure collection actually exists and we have the correct ID
	collection, err = client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading Scheduler Job Collection %q after create/update (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*collection.ID)

	return resourceArmSchedulerJobCollectionRead(d, meta)
}

func resourceArmSchedulerJobCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).scheduler.JobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobCollections"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Reading Scheduler Job Collection %q (resource group %q)", name, resourceGroup)

	collection, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(collection.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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

func resourceArmSchedulerJobCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).scheduler.JobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobCollections"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Deleting Scheduler Job Collection %q (resource group %q)", name, resourceGroup)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error issuing delete request for Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandAzureArmSchedulerJobCollectionQuota(d *schema.ResourceData) *scheduler.JobCollectionQuota {
	if qb, ok := d.Get("quota").([]interface{}); ok && len(qb) > 0 {
		quota := scheduler.JobCollectionQuota{
			MaxRecurrence: &scheduler.JobMaxRecurrence{},
		}

		quotaBlock := qb[0].(map[string]interface{})

		if v, ok := quotaBlock["max_job_count"].(int); ok {
			quota.MaxJobCount = utils.Int32(int32(v))
		}
		if v, ok := quotaBlock["max_recurrence_frequency"].(string); ok {
			quota.MaxRecurrence.Frequency = scheduler.RecurrenceFrequency(v)
		}
		if v, ok := quotaBlock["max_recurrence_interval"].(int); ok && v > 0 {
			quota.MaxRecurrence.Interval = utils.Int32(int32(v))
		} else if v, ok := quotaBlock["max_retry_interval"].(int); ok && v > 0 { //todo remove once max_retry_interval is removed
			quota.MaxRecurrence.Interval = utils.Int32(int32(v))
		}

		return &quota
	}

	return nil
}

func flattenAzureArmSchedulerJobCollectionQuota(quota *scheduler.JobCollectionQuota) []interface{} {

	if quota == nil {
		return nil
	}

	quotaBlock := make(map[string]interface{})

	if v := quota.MaxJobCount; v != nil {
		quotaBlock["max_job_count"] = *v
	}
	if recurrence := quota.MaxRecurrence; recurrence != nil {
		if v := recurrence.Interval; v != nil {
			quotaBlock["max_recurrence_interval"] = *v
			quotaBlock["max_retry_interval"] = *v //todo remove once max_retry_interval is retired
		}

		quotaBlock["max_recurrence_frequency"] = string(recurrence.Frequency)
	}

	return []interface{}{quotaBlock}
}
