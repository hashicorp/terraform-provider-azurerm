package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryTriggerBlobEvent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryTriggerBlobEventCreateUpdate,
		Read:   resourceDataFactoryTriggerBlobEventRead,
		Update: resourceDataFactoryTriggerBlobEventCreateUpdate,
		Delete: resourceDataFactoryTriggerBlobEventDelete,

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
				ValidateFunc: validate.DataFactoryID,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"events": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Microsoft.Storage.BlobCreated",
						"Microsoft.Storage.BlobDeleted",
					}, false),
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

			"blob_path_begins_with": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"blob_path_begins_with", "blob_path_ends_with"},
			},

			"blob_path_ends_with": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"blob_path_begins_with", "blob_path_ends_with"},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"ignore_empty_blobs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDataFactoryTriggerBlobEventCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := parse.DataFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_blob_event", id.ID())
		}
	}

	blobEventProps := &datafactory.BlobEventsTrigger{
		BlobEventsTriggerTypeProperties: &datafactory.BlobEventsTriggerTypeProperties{
			IgnoreEmptyBlobs: utils.Bool(d.Get("ignore_empty_blobs").(bool)),
			Events:           expandDataFactoryTriggerBlobEvents(d.Get("events").(*pluginsdk.Set).List()),
			Scope:            utils.String(d.Get("storage_account_id").(string)),
		},
		Description: utils.String(d.Get("description").(string)),
		Pipelines:   expandDataFactoryTriggerPipeline(d.Get("pipeline").(*pluginsdk.Set).List()),
		Type:        datafactory.TypeBasicTriggerTypeBlobEventsTrigger,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		blobEventProps.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		blobEventProps.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("blob_path_begins_with"); ok {
		blobEventProps.BlobPathBeginsWith = utils.String(v.(string))
	}

	if v, ok := d.GetOk("blob_path_ends_with"); ok {
		blobEventProps.BlobPathEndsWith = utils.String(v.(string))
	}

	trigger := datafactory.TriggerResource{
		Properties: blobEventProps,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, trigger, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryTriggerBlobEventRead(d, meta)
}

func resourceDataFactoryTriggerBlobEventRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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

	blobEventsTrigger, ok := resp.Properties.AsBlobEventsTrigger()
	if !ok {
		return fmt.Errorf("classifiying %s: Expected: %q", id, datafactory.TypeBasicTriggerTypeBlobEventsTrigger)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", parse.NewDataFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	d.Set("additional_properties", blobEventsTrigger.AdditionalProperties)
	d.Set("description", blobEventsTrigger.Description)

	if err := d.Set("annotations", flattenDataFactoryAnnotations(blobEventsTrigger.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("pipeline", flattenDataFactoryTriggerPipeline(blobEventsTrigger.Pipelines)); err != nil {
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

func resourceDataFactoryTriggerBlobEventDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TriggerID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDataFactoryTriggerBlobEvents(input []interface{}) *[]datafactory.BlobEventTypes {
	result := make([]datafactory.BlobEventTypes, 0)
	for _, item := range input {
		result = append(result, datafactory.BlobEventTypes(item.(string)))
	}
	return &result
}

func expandDataFactoryTriggerPipeline(input []interface{}) *[]datafactory.TriggerPipelineReference {
	if len(input) == 0 {
		return nil
	}

	result := make([]datafactory.TriggerPipelineReference, 0)
	for _, item := range input {
		raw := item.(map[string]interface{})

		// issue https://github.com/hashicorp/terraform-plugin-sdk/issues/588
		// once it's resolved, we could remove the check empty logic
		name := raw["name"].(string)
		if name == "" {
			continue
		}

		result = append(result, datafactory.TriggerPipelineReference{
			PipelineReference: &datafactory.PipelineReference{
				ReferenceName: utils.String(raw["name"].(string)),
				Type:          utils.String("PipelineReference"),
			},
			Parameters: raw["parameters"].(map[string]interface{}),
		})
	}
	return &result
}

func flattenDataFactoryTriggerBlobEvents(input *[]datafactory.BlobEventTypes) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		result = append(result, string(item))
	}
	return result
}

func flattenDataFactoryTriggerPipeline(input *[]datafactory.TriggerPipelineReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		name := ""
		if item.PipelineReference != nil && item.PipelineReference.ReferenceName != nil {
			name = *item.PipelineReference.ReferenceName
		}

		result = append(result, map[string]interface{}{
			"name":       name,
			"parameters": item.Parameters,
		})
	}
	return result
}
