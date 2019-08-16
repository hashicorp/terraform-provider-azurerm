package azurerm

import (
	"fmt"
	"log"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubCreateUpdate,
		Read:   resourceArmEventHubRead,
		Update: resourceArmEventHubCreateUpdate,
		Delete: resourceArmEventHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			// TODO: remove me in the next major version
			"location": azure.SchemaLocationDeprecated(),

			"partition_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateEventHubPartitionCount,
			},

			"message_retention": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateEventHubMessageRetentionCount,
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
										ValidateFunc: validateEventHubArchiveNameFormat,
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

func resourceArmEventHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.EventHubsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM EventHub creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub", *existing.ID)
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
		captureDescription, err := expandEventHubCaptureDescription(d)
		if err != nil {
			return fmt.Errorf("Error expanding EventHub Capture Description: %s", err)
		}

		parameters.Properties.CaptureDescription = captureDescription
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubRead(d, meta)
}

func resourceArmEventHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.EventHubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["eventhubs"]
	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure EventHub %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

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

func resourceArmEventHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.EventHubsClient
	ctx := meta.(*ArmClient).StopContext
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["eventhubs"]
	resp, err := client.Delete(ctx, resourceGroup, namespaceName, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing delete request for EventHub %q (resource group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func validateEventHubPartitionCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)

	if !(32 >= value && value >= 1) {
		errors = append(errors, fmt.Errorf("EventHub Partition Count has to be between 1 and 32"))
	}

	return warnings, errors
}

func validateEventHubMessageRetentionCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)

	if !(7 >= value && value >= 1) {
		errors = append(errors, fmt.Errorf("EventHub Retention Count has to be between 1 and 7"))
	}

	return warnings, errors
}

func validateEventHubArchiveNameFormat(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	requiredComponents := []string{
		"{Namespace}",
		"{EventHub}",
		"{PartitionId}",
		"{Year}",
		"{Month}",
		"{Day}",
		"{Hour}",
		"{Minute}",
		"{Second}",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(value, component) {
			errors = append(errors, fmt.Errorf("%s needs to contain %q", k, component))
		}
	}

	return warnings, errors
}

func expandEventHubCaptureDescription(d *schema.ResourceData) (*eventhub.CaptureDescription, error) {
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

	return &captureDescription, nil
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
