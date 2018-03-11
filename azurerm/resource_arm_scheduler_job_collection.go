package azurerm

import (
	"fmt"
	"log"

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

			"tags": tagsSchema(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
				Default:          "enabled",
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(scheduler.Enabled),
					string(scheduler.Suspended),
					string(scheduler.Disabled),
					string(scheduler.Deleted),
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(scheduler.Minute),
								string(scheduler.Hour),
								string(scheduler.Day),
								string(scheduler.Week),
								string(scheduler.Month),
							}, true),
						},

						//should this be max_retry_interval ? given that is what the documentation implies
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
	client := meta.(*ArmClient).schedulerJobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating/updating Scheduler Job Collection %q (resource group %q)", name, resourceGroup)

	sku := scheduler.Sku{
		Name: scheduler.SkuDefinition(d.Get("sku").(string)),
	}

	properties := scheduler.JobCollectionProperties{
		Sku: &sku,
	}
	if state, ok := d.Get("state").(string); ok {
		properties.State = scheduler.JobCollectionState(state)
	}

	if qb, ok := d.Get("quota").([]interface{}); ok && len(qb) > 0 {
		recurrence := scheduler.JobMaxRecurrence{}
		quota := scheduler.JobCollectionQuota{
			MaxRecurrence: &recurrence,
		}

		quotaBlock := qb[0].(map[string]interface{})

		if v, ok := quotaBlock["max_job_count"].(int); ok {
			quota.MaxJobCount = utils.Int32(int32(v))
		}

		if v, ok := quotaBlock["max_recurrence_frequency"].(string); ok {
			recurrence.Frequency = scheduler.RecurrenceFrequency(v)
		}
		if v, ok := quotaBlock["max_recurrence_interval"].(int); ok {
			recurrence.Interval = utils.Int32(int32(v))
		}

		properties.Quota = &quota
	}

	collection := scheduler.JobCollectionDefinition{
		Location:   utils.String(location),
		Tags:       expandTags(tags),
		Properties: &properties,
	}

	//create job collection
	collection, err := client.CreateOrUpdate(ctx, resourceGroup, name, collection)
	if err != nil {
		return fmt.Errorf("Error creating/updating Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	//ensure collection actually exists and we have the correct ID
	collection, err = client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*collection.ID)

	return resourceArmSchedulerJobCollectionPopulate(d, resourceGroup, &collection)
}

func resourceArmSchedulerJobCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobCollections"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Reading Scheduler Job Collection %q (resource group %q)", name, resourceGroup)

	collection, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		//TODO why is this in utils not response?
		if utils.ResponseWasNotFound(collection.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Scheduler Job Collection Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmSchedulerJobCollectionPopulate(d, resourceGroup, &collection)
}

func resourceArmSchedulerJobCollectionPopulate(d *schema.ResourceData, resourceGroup string, collection *scheduler.JobCollectionDefinition) error {

	//standard properties
	d.Set("name", collection.Name)
	d.Set("location", azureRMNormalizeLocation(*collection.Location))
	d.Set("resource_group_name", resourceGroup)
	flattenAndSetTags(d, collection.Tags)

	//resource specific
	if properties := collection.Properties; properties != nil {
		if sku := properties.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}
		d.Set("state", string(properties.State))

		if quota := properties.Quota; quota != nil {
			quotaBlock := make(map[string]interface{})

			if v := quota.MaxJobCount; v != nil {
				quotaBlock["max_job_count"] = *v
			}

			if recurrence := quota.MaxRecurrence; recurrence != nil {
				if v := recurrence.Interval; v != nil {
					quotaBlock["max_recurrence_interval"] = *v
				}

				quotaBlock["max_recurrence_frequency"] = string(recurrence.Frequency)
			}

			d.Set("quota", []interface{}{quotaBlock})
		}
	}

	return nil
}

func resourceArmSchedulerJobCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).schedulerJobCollectionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Scheduler Job Collection %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
