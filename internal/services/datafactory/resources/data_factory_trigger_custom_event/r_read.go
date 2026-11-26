// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_custom_event

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryTriggerCustomEventRead(d *pluginsdk.ResourceData, meta any) error {
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
		if response.WasNotFound(resp.Response.Response) {
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

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(trigger.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("pipeline", helpers.FlattenDataFactoryTriggerPipeline(trigger.Pipelines)); err != nil {
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
