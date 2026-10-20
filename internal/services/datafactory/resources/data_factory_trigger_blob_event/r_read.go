// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_blob_event

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

func resourceDataFactoryTriggerBlobEventRead(d *pluginsdk.ResourceData, meta any) error {
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

	blobEventsTrigger, ok := resp.Properties.AsBlobEventsTrigger()
	if !ok {
		return fmt.Errorf("classifying %s: Expected: %q", id, datafactory.TypeBasicTriggerTypeBlobEventsTrigger)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	d.Set("activated", blobEventsTrigger.RuntimeState == datafactory.TriggerRuntimeStateStarted)
	d.Set("additional_properties", blobEventsTrigger.AdditionalProperties)
	d.Set("description", blobEventsTrigger.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(blobEventsTrigger.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("pipeline", helpers.FlattenDataFactoryTriggerPipeline(blobEventsTrigger.Pipelines)); err != nil {
		return fmt.Errorf("setting `pipeline`: %+v", err)
	}

	if props := blobEventsTrigger.BlobEventsTriggerTypeProperties; props != nil {
		d.Set("storage_account_id", props.Scope)
		d.Set("blob_path_begins_with", props.BlobPathBeginsWith)
		d.Set("blob_path_ends_with", props.BlobPathEndsWith)
		d.Set("ignore_empty_blobs", props.IgnoreEmptyBlobs)

		if err := d.Set("events", flattenDataFactoryTriggerBlobEvents(props.Events)); err != nil {
			return fmt.Errorf("setting `events`: %+v", err)
		}
	}

	return nil
}

func flattenDataFactoryTriggerBlobEvents(input *[]datafactory.BlobEventTypes) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, item := range *input {
		result = append(result, string(item))
	}
	return result
}
