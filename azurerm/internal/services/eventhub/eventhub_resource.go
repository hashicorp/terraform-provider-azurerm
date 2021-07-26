package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var eventHubResourceName = "azurerm_eventhub"

func resourceEventHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHubCreateUpdate,
		Read:   resourceEventHubRead,
		Update: resourceEventHubCreateUpdate,
		Delete: resourceEventHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EventHubID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"partition_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubPartitionCount,
			},

			"message_retention": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubMessageRetentionCount,
			},

			"capture_description": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"skip_empty_archives": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"encoding": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.Avro),
								string(eventhub.AvroDeflate),
							}, true),
						},
						"interval_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 900),
						},
						"size_limit_in_bytes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      314572800,
							ValidateFunc: validation.IntBetween(10485760, 524288000),
						},
						"destination": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"EventHubArchive.AzureBlockBlob",
											// TODO: support `EventHubArchive.AzureDataLake` once supported in the Swagger / SDK
											// https://github.com/Azure/azure-rest-api-specs/issues/2255
											// BlobContainerName & StorageAccountID can then become Optional
										}, false),
									},
									"archive_name_format": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ValidateEventHubArchiveNameFormat,
									},
									"blob_container_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"storage_account_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
							},
						},
					},
				},
			},

			"partition_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
		},
	}
}

func resourceEventHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM EventHub creation.")

	id := parse.NewEventHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_eventhub", id.ID())
		}
	}

	partitionCount := int64(d.Get("partition_count").(int))
	messageRetention := int64(d.Get("message_retention").(int))

	parameters := eventhub.Model{
		Properties: &eventhub.Properties{
			PartitionCount:         &partitionCount,
			MessageRetentionInDays: &messageRetention,
		},
	}

	if _, ok := d.GetOk("capture_description"); ok {
		parameters.Properties.CaptureDescription = expandEventHubCaptureDescription(d)
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.NamespaceName, id.Name, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceEventHubRead(d, meta)
}

func resourceEventHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("partition_count", props.PartitionCount)
		d.Set("message_retention", props.MessageRetentionInDays)
		d.Set("partition_ids", props.PartitionIds)

		captureDescription := flattenEventHubCaptureDescription(props.CaptureDescription)
		if err := d.Set("capture_description", captureDescription); err != nil {
			return err
		}
	}

	return nil
}

func resourceEventHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandEventHubCaptureDescription(d *schema.ResourceData) *eventhub.CaptureDescription {
	inputs := d.Get("capture_description").([]interface{})
	input := inputs[0].(map[string]interface{})

	enabled := input["enabled"].(bool)
	encoding := input["encoding"].(string)
	intervalInSeconds := input["interval_in_seconds"].(int)
	sizeLimitInBytes := input["size_limit_in_bytes"].(int)
	skipEmptyArchives := input["skip_empty_archives"].(bool)

	captureDescription := eventhub.CaptureDescription{
		Enabled:           utils.Bool(enabled),
		Encoding:          eventhub.EncodingCaptureDescription(encoding),
		IntervalInSeconds: utils.Int32(int32(intervalInSeconds)),
		SizeLimitInBytes:  utils.Int32(int32(sizeLimitInBytes)),
		SkipEmptyArchives: utils.Bool(skipEmptyArchives),
	}

	if v, ok := input["destination"]; ok {
		destinations := v.([]interface{})
		if len(destinations) > 0 {
			destination := destinations[0].(map[string]interface{})

			destinationName := destination["name"].(string)
			archiveNameFormat := destination["archive_name_format"].(string)
			blobContainerName := destination["blob_container_name"].(string)
			storageAccountId := destination["storage_account_id"].(string)

			captureDescription.Destination = &eventhub.Destination{
				Name: utils.String(destinationName),
				DestinationProperties: &eventhub.DestinationProperties{
					ArchiveNameFormat:        utils.String(archiveNameFormat),
					BlobContainer:            utils.String(blobContainerName),
					StorageAccountResourceID: utils.String(storageAccountId),
				},
			}
		}
	}

	return &captureDescription
}

func flattenEventHubCaptureDescription(description *eventhub.CaptureDescription) []interface{} {
	results := make([]interface{}, 0)

	if description != nil {
		output := make(map[string]interface{})

		if enabled := description.Enabled; enabled != nil {
			output["enabled"] = *enabled
		}

		if skipEmptyArchives := description.SkipEmptyArchives; skipEmptyArchives != nil {
			output["skip_empty_archives"] = *skipEmptyArchives
		}

		output["encoding"] = string(description.Encoding)

		if interval := description.IntervalInSeconds; interval != nil {
			output["interval_in_seconds"] = *interval
		}

		if size := description.SizeLimitInBytes; size != nil {
			output["size_limit_in_bytes"] = *size
		}

		if destination := description.Destination; destination != nil {
			destinationOutput := make(map[string]interface{})

			if name := destination.Name; name != nil {
				destinationOutput["name"] = *name
			}

			if props := destination.DestinationProperties; props != nil {
				if archiveNameFormat := props.ArchiveNameFormat; archiveNameFormat != nil {
					destinationOutput["archive_name_format"] = *archiveNameFormat
				}
				if blobContainerName := props.BlobContainer; blobContainerName != nil {
					destinationOutput["blob_container_name"] = *blobContainerName
				}
				if storageAccountId := props.StorageAccountResourceID; storageAccountId != nil {
					destinationOutput["storage_account_id"] = *storageAccountId
				}
			}

			output["destination"] = []interface{}{destinationOutput}
		}

		results = append(results, output)
	}

	return results
}
