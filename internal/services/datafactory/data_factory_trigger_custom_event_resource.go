// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryTriggerCustomEvent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryTriggerCustomEventCreateUpdate,
		Read:   resourceDataFactoryTriggerCustomEventRead,
		Update: resourceDataFactoryTriggerCustomEventCreateUpdate,
		Delete: resourceDataFactoryTriggerCustomEventDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.TriggerID(id)
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

			"eventgrid_topic_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: topics.ValidateTopicID,
			},

			"events": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"pipeline": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
						},

						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"activated": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subject_begins_with": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"subject_begins_with", "subject_ends_with"},
			},

			"subject_ends_with": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"subject_begins_with", "subject_ends_with"},
			},
		},
	}
}

func resourceDataFactoryTriggerCustomEventCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_custom_event", id.ID())
		}
	} else {
		future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to stop %s: %+v", id, err)
		}
	}

	events := d.Get("events").(*pluginsdk.Set).List()
	trigger := &datafactory.CustomEventsTrigger{
		CustomEventsTriggerTypeProperties: &datafactory.CustomEventsTriggerTypeProperties{
			Events: &events,
			Scope:  utils.String(d.Get("eventgrid_topic_id").(string)),
		},
		Description: utils.String(d.Get("description").(string)),
		Pipelines:   expandDataFactoryTriggerPipeline(d.Get("pipeline").(*pluginsdk.Set).List()),
		Type:        datafactory.TypeBasicTriggerTypeCustomEventsTrigger,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		trigger.Annotations = &annotations
	}

	if v, ok := d.GetOk("subject_begins_with"); ok {
		trigger.SubjectBeginsWith = utils.String(v.(string))
	}

	if v, ok := d.GetOk("subject_ends_with"); ok {
		trigger.SubjectEndsWith = utils.String(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		trigger.AdditionalProperties = v.(map[string]interface{})
	}

	resource := datafactory.TriggerResource{
		Properties: trigger,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, resource, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.Get("activated").(bool) {
		future, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on start %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceDataFactoryTriggerCustomEventRead(d, meta)
}

func resourceDataFactoryTriggerCustomEventRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	trigger, ok := resp.Properties.AsCustomEventsTrigger()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q", id, datafactory.TypeBasicTriggerTypeCustomEventsTrigger)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	d.Set("activated", trigger.RuntimeState == datafactory.TriggerRuntimeStateStarted)
	d.Set("additional_properties", trigger.AdditionalProperties)
	d.Set("description", trigger.Description)

	if err := d.Set("annotations", flattenDataFactoryAnnotations(trigger.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("pipeline", flattenDataFactoryTriggerPipeline(trigger.Pipelines)); err != nil {
		return fmt.Errorf("setting `pipeline`: %+v", err)
	}

	if props := trigger.CustomEventsTriggerTypeProperties; props != nil {
		d.Set("eventgrid_topic_id", props.Scope)
		d.Set("events", props.Events)
		d.Set("subject_begins_with", props.SubjectBeginsWith)
		d.Set("subject_ends_with", props.SubjectEndsWith)
	}

	return nil
}

func resourceDataFactoryTriggerCustomEventDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		return fmt.Errorf("stopping %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to stop %s: %+v", id, err)
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
