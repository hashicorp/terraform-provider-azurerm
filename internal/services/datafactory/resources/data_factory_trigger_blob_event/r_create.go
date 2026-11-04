// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_blob_event

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

func resourceDataFactoryTriggerBlobEventCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
				return tf.ImportAsExistsError("azurerm_data_factory_trigger_blob_event", id.ID())
			}
		}
	}
	if !d.IsNewResource() {
		future, err := client.Stop(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to stop %s: %+v", id, err)
		}
	}

	blobEventProps := &datafactory.BlobEventsTrigger{
		BlobEventsTriggerTypeProperties: &datafactory.BlobEventsTriggerTypeProperties{
			IgnoreEmptyBlobs: pointer.To(d.Get("ignore_empty_blobs").(bool)),
			Events:           expandDataFactoryTriggerBlobEvents(d.Get("events").(*pluginsdk.Set).List()),
			Scope:            pointer.To(d.Get("storage_account_id").(string)),
		},
		Description: pointer.To(d.Get("description").(string)),
		Pipelines:   helpers.ExpandDataFactoryTriggerPipeline(d.Get("pipeline").(*pluginsdk.Set).List()),
		Type:        datafactory.TypeBasicTriggerTypeBlobEventsTrigger,
	}

	if v, ok := d.GetOk("annotations"); ok {
		blobEventProps.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		blobEventProps.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("blob_path_begins_with"); ok {
		blobEventProps.BlobPathBeginsWith = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("blob_path_ends_with"); ok {
		blobEventProps.BlobPathEndsWith = pointer.To(v.(string))
	}

	trigger := datafactory.TriggerResource{
		Properties: blobEventProps,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, trigger, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if v, ok := d.GetOk("activated"); ok && v.(bool) {
		future, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting on start %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceDataFactoryTriggerBlobEventRead(d, meta)
}

func expandDataFactoryTriggerBlobEvents(input []any) *[]datafactory.BlobEventTypes {
	result := make([]datafactory.BlobEventTypes, 0)
	for _, item := range input {
		result = append(result, datafactory.BlobEventTypes(item.(string)))
	}
	return &result
}
