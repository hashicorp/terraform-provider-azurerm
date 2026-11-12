// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_tumbling_window

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerTumblingWindowCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_trigger_tumbling_window", id.ID())
			}
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

	startTime, err := time.Parse(time.RFC3339, d.Get("start_time").(string))
	if err != nil {
		return err
	}

	props := &datafactory.TumblingWindowTrigger{
		TumblingWindowTriggerTypeProperties: &datafactory.TumblingWindowTriggerTypeProperties{
			Frequency:      datafactory.TumblingWindowFrequency(d.Get("frequency").(string)),
			Interval:       pointer.To(int32(d.Get("interval").(int))),
			MaxConcurrency: pointer.To(int32(d.Get("max_concurrency").(int))),
			RetryPolicy:    expandDataFactoryTriggerTumblingWindowRetryPolicy(d.Get("retry").([]any)),
			DependsOn:      expandDataFactoryTriggerDependency(d.Get("trigger_dependency").(*pluginsdk.Set).List()),
			StartTime:      &date.Time{Time: startTime},
		},
		Description: pointer.To(d.Get("description").(string)),
		Pipeline:    expandDataFactoryTriggerSinglePipeline(d.Get("pipeline").([]any)),
		Type:        datafactory.TypeBasicTriggerTypeTumblingWindowTrigger,
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return err
		}
		props.EndTime = &date.Time{Time: t}
	}

	if v, ok := d.GetOk("delay"); ok {
		props.Delay = v.(string)
	}

	if v, ok := d.GetOk("annotations"); ok {
		props.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		props.AdditionalProperties = v.(map[string]any)
	}

	trigger := datafactory.TriggerResource{
		Properties: props,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, trigger, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.Get("activated").(bool) {
		future, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on start %s: %+v", id, err)
		}
	}

	return resourceDataFactoryTriggerTumblingWindowRead(d, meta)
}
