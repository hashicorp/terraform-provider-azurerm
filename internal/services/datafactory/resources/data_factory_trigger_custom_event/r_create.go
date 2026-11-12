// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_custom_event

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerCustomEventCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
				return tf.ImportAsExistsError("azurerm_data_factory_trigger_custom_event", id.ID())
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

	events := d.Get("events").(*pluginsdk.Set).List()
	trigger := &datafactory.CustomEventsTrigger{
		CustomEventsTriggerTypeProperties: &datafactory.CustomEventsTriggerTypeProperties{
			Events: &events,
			Scope:  pointer.To(d.Get("eventgrid_topic_id").(string)),
		},
		Description: pointer.To(d.Get("description").(string)),
		Pipelines:   helpers.ExpandDataFactoryTriggerPipeline(d.Get("pipeline").(*pluginsdk.Set).List()),
		Type:        datafactory.TypeBasicTriggerTypeCustomEventsTrigger,
	}

	if v, ok := d.GetOk("annotations"); ok {
		trigger.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("subject_begins_with"); ok {
		trigger.SubjectBeginsWith = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("subject_ends_with"); ok {
		trigger.SubjectEndsWith = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		trigger.AdditionalProperties = v.(map[string]any)
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
