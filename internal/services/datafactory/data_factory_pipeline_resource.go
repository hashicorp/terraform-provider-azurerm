// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryPipeline() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryPipelineCreateUpdate,
		Read:   resourceDataFactoryPipelineRead,
		Update: resourceDataFactoryPipelineCreateUpdate,
		Delete: resourceDataFactoryPipelineDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PipelineID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"variables": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"activities_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"concurrency": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 50),
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"moniter_metrics_after_duration": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDataFactoryPipelineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	subscriptionId := meta.(*clients.Client).DataFactory.PipelinesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewPipelineID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_pipeline", id.ID())
		}
	}

	pipeline := &datafactory.Pipeline{
		Parameters:  expandDataFactoryParameters(d.Get("parameters").(map[string]interface{})),
		Variables:   expandDataFactoryVariables(d.Get("variables").(map[string]interface{})),
		Description: utils.String(d.Get("description").(string)),
	}

	if v, ok := d.GetOk("activities_json"); ok {
		activities, err := deserializeDataFactoryPipelineActivities(v.(string))
		if err != nil {
			return fmt.Errorf("parsing 'activities_json' for Data Factory %s: %+v", id, err)
		}
		pipeline.Activities = activities
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		pipeline.Annotations = &annotations
	} else {
		annotations := make([]interface{}, 0)
		pipeline.Annotations = &annotations
	}

	if v, ok := d.GetOk("concurrency"); ok {
		pipeline.Concurrency = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("moniter_metrics_after_duration"); ok {
		pipeline.Policy = &datafactory.PipelinePolicy{
			ElapsedTimeMetric: &datafactory.PipelineElapsedTimeMetricPolicy{
				Duration: v.(string),
			},
		}
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		pipeline.Folder = &datafactory.PipelineFolder{
			Name: &name,
		}
	}

	config := datafactory.PipelineResource{
		Pipeline: pipeline,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, config, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryPipelineRead(d, meta)
}

func resourceDataFactoryPipelineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PipelineID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Data Factory Pipeline %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("reading the state of Data Factory Pipeline %q: %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	if props := resp.Pipeline; props != nil {
		d.Set("description", props.Description)

		parameters := flattenDataFactoryParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		annotations := flattenDataFactoryAnnotations(props.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		concurrency := 0
		if props.Concurrency != nil {
			concurrency = int(*props.Concurrency)
		}
		d.Set("concurrency", concurrency)

		elapsedTimeMetricDuration := ""
		if props.Policy != nil && props.Policy.ElapsedTimeMetric != nil && props.Policy.ElapsedTimeMetric.Duration != nil {
			if v, ok := props.Policy.ElapsedTimeMetric.Duration.(string); ok {
				elapsedTimeMetricDuration = v
			}
		}
		d.Set("moniter_metrics_after_duration", elapsedTimeMetricDuration)

		if folder := props.Folder; folder != nil {
			if folder.Name != nil {
				d.Set("folder", folder.Name)
			}
		}

		variables := flattenDataFactoryVariables(props.Variables)
		if err := d.Set("variables", variables); err != nil {
			return fmt.Errorf("setting `variables`: %+v", err)
		}

		if activities := props.Activities; activities != nil {
			activitiesJson, err := serializeDataFactoryPipelineActivities(activities)
			if err != nil {
				return fmt.Errorf("serializing `activities_json`: %+v", err)
			}
			if err := d.Set("activities_json", activitiesJson); err != nil {
				return fmt.Errorf("setting `activities_json`: %+v", err)
			}
		}
	}

	return nil
}

func resourceDataFactoryPipelineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PipelineID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
