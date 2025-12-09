// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var eventHubResourceName = "azurerm_eventhub"

func resourceEventHub() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceEventHubCreate,
		Read:   resourceEventHubRead,
		Update: resourceEventHubUpdate,
		Delete: resourceEventHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubs.ParseEventhubID(id)
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
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"partition_count": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.ValidateEventHubPartitionCount,
			},

			"message_retention": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateEventHubMessageRetentionCount,
				ExactlyOneOf: []string{"retention_description", "message_retention"},
			},

			"retention_description": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"cleanup_policy": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhubs.CleanupPolicyRetentionDescriptionDelete),
								string(eventhubs.CleanupPolicyRetentionDescriptionCompact),
							}, false),
						},

						"retention_time_in_hours": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"retention_description.0.tombstone_retention_time_in_hours", "retention_description.0.retention_time_in_hours"},
						},

						"tombstone_retention_time_in_hours": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"retention_description.0.retention_time_in_hours", "retention_description.0.tombstone_retention_time_in_hours"},
						},
					},
				},
			},

			"capture_description": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"skip_empty_archives": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"encoding": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhubs.EncodingCaptureDescriptionAvro),
								string(eventhubs.EncodingCaptureDescriptionAvroDeflate),
							}, false),
						},
						"interval_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      300,
							ValidateFunc: validation.IntBetween(60, 900),
						},
						"size_limit_in_bytes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      314572800,
							ValidateFunc: validation.IntBetween(10485760, 524288000),
						},
						"destination": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"EventHubArchive.AzureBlockBlob",
											// TODO: support `EventHubArchive.AzureDataLake` once supported in the Swagger / SDK
											// https://github.com/Azure/azure-rest-api-specs/issues/2255
											// BlobContainerName & StorageAccountID can then become Optional
										}, false),
									},
									"archive_name_format": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validate.ValidateEventHubArchiveNameFormat,
									},
									"blob_container_name": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},
									"storage_account_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: commonids.ValidateStorageAccountID,
									},
									"storage_authentication": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(eventhubs.CaptureIdentityTypeSystemAssigned),
											string(eventhubs.CaptureIdentityTypeUserAssigned),
										}, false),
									},
									"storage_authentication_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: commonids.ValidateUserAssignedIdentityID,
									},
								},
							},
						},
					},
				},
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(eventhubs.EntityStatusActive),
				ValidateFunc: validation.StringInSlice([]string{
					string(eventhubs.EntityStatusActive),
					string(eventhubs.EntityStatusDisabled),
					string(eventhubs.EntityStatusSendDisabled),
				}, false),
			},

			"partition_ids": {
				Type:     pluginsdk.TypeSet,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["namespace_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"namespace_id", "namespace_name"},
			ValidateFunc: namespaces.ValidateNamespaceID,
		}

		r.Schema["namespace_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ValidateEventHubNamespaceName(),
			ExactlyOneOf: []string{"namespace_id", "namespace_name"},
			Deprecated:   "`namespace_name` has been deprecated in favour of `namespace_id` and will be removed in v5.0 of the AzureRM Provider",
		}

		r.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"namespace_id", "resource_group_name"},
			ValidateFunc: resourcegroups.ValidateName,
			Deprecated:   "`resource_group_name` has been deprecated in favour of `namespace_id` and will be removed in v5.0 of the AzureRM Provider",
		}
	}

	return r
}

func resourceEventHubCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM EventHub creation.")

	namespaceName := ""
	resourceGroupName := ""
	if v := d.Get("namespace_id").(string); v != "" {
		namespaceId, err := namespaces.ParseNamespaceID(v)
		if err != nil {
			return err
		}
		namespaceName = namespaceId.NamespaceName
		resourceGroupName = namespaceId.ResourceGroupName
	}

	if !features.FivePointOh() && namespaceName == "" {
		namespaceName = d.Get("namespace_name").(string)
		resourceGroupName = d.Get("resource_group_name").(string)
	}

	id := eventhubs.NewEventhubID(subscriptionId, resourceGroupName, namespaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_eventhub", id.ID())
		}
	}

	eventhubStatus := eventhubs.EntityStatus(d.Get("status").(string))
	parameters := eventhubs.Eventhub{
		Properties: &eventhubs.EventhubProperties{
			PartitionCount: utils.Int64(int64(d.Get("partition_count").(int))),
			Status:         &eventhubStatus,
		},
	}

	if _, ok := d.GetOk("retention_description"); ok {
		parameters.Properties.RetentionDescription = expandEventHubRetentionDescription(d)
	}

	if _, ok := d.GetOk("message_retention"); ok {
		parameters.Properties.MessageRetentionInDays = pointer.To(int64(d.Get("message_retention").(int)))
	}

	if _, ok := d.GetOk("capture_description"); ok {
		parameters.Properties.CaptureDescription = expandEventHubCaptureDescription(d)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceEventHubRead(d, meta)
}

func resourceEventHubUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM EventHub update.")

	id, err := eventhubs.ParseEventhubID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("partition_count") {
		o, n := d.GetChange("partition_count")
		if o.(int) > n.(int) {
			return fmt.Errorf("`partition_count` cannot be decreased")
		}

		client := meta.(*clients.Client).Eventhub.NamespacesClient
		namespaceId := namespaces.NewNamespaceID(subscriptionId, id.ResourceGroupName, id.NamespaceName)
		resp, err := client.Get(ctx, namespaceId)
		if err != nil {
			return err
		}
		if model := resp.Model; model != nil {
			if model.Sku.Name != namespaces.SkuNamePremium {
				return fmt.Errorf("`partition_count` cannot be changed unless the namespace sku is `Premium`")
			}
		}
	}

	eventhubStatus := eventhubs.EntityStatus(d.Get("status").(string))
	parameters := eventhubs.Eventhub{
		Properties: &eventhubs.EventhubProperties{
			PartitionCount:         utils.Int64(int64(d.Get("partition_count").(int))),
			Status:                 &eventhubStatus,
			MessageRetentionInDays: utils.Int64(int64(d.Get("message_retention").(int))),
			CaptureDescription:     expandEventHubCaptureDescription(d),
		},
	}

	if d.HasChange("capture_description") {
		parameters.Properties.CaptureDescription = expandEventHubCaptureDescription(d)
	}

	if d.HasChange("retention_description") {
		parameters.Properties.RetentionDescription = expandEventHubRetentionDescription(d)
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceEventHubRead(d, meta)
}

func resourceEventHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubs.ParseEventhubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.EventhubName)

	if !features.FivePointOh() {
		d.Set("namespace_name", id.NamespaceName)
		d.Set("resource_group_name", id.ResourceGroupName)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	d.Set("namespace_id", namespaceId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("partition_count", props.PartitionCount)
			d.Set("partition_ids", props.PartitionIds)
			d.Set("status", string(*props.Status))

			captureDescription := flattenEventHubCaptureDescription(props.CaptureDescription)
			if err := d.Set("capture_description", captureDescription); err != nil {
				return err
			}

			// TODO - the `props.RetentionDescription.TombstoneRetentionTimeInHours == nil` check can be removed when https://github.com/Azure/azure-rest-api-specs/issues/36018 is fixed
			if props.RetentionDescription == nil || props.RetentionDescription.TombstoneRetentionTimeInHours == nil {
				d.Set("message_retention", props.MessageRetentionInDays)
			}
			retentionDescription := flattenEventHubRetentionDescription(props.RetentionDescription)
			if err := d.Set("retention_description", retentionDescription); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceEventHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubs.ParseEventhubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandEventHubRetentionDescription(d *pluginsdk.ResourceData) *eventhubs.RetentionDescription {
	inputs := d.Get("retention_description").([]interface{})
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	cleanupPolicy := input["cleanup_policy"].(string)

	retentionDescription := eventhubs.RetentionDescription{
		CleanupPolicy: pointer.To(eventhubs.CleanupPolicyRetentionDescription(cleanupPolicy)),
	}

	if cleanupPolicy == string(eventhubs.CleanupPolicyRetentionDescriptionDelete) {
		retentionTimeInHours := input["retention_time_in_hours"].(int)
		retentionDescription.RetentionTimeInHours = pointer.FromInt64(int64(retentionTimeInHours))
	} else {
		tombstoneRetentionTimeInHours := input["tombstone_retention_time_in_hours"].(int)
		retentionDescription.TombstoneRetentionTimeInHours = pointer.FromInt64(int64(tombstoneRetentionTimeInHours))
	}

	return pointer.To(retentionDescription)
}

func expandEventHubCaptureDescription(d *pluginsdk.ResourceData) *eventhubs.CaptureDescription {
	inputs := d.Get("capture_description").([]interface{})
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}
	input := inputs[0].(map[string]interface{})

	enabled := input["enabled"].(bool)
	encoding := input["encoding"].(string)
	intervalInSeconds := input["interval_in_seconds"].(int)
	sizeLimitInBytes := input["size_limit_in_bytes"].(int)
	skipEmptyArchives := input["skip_empty_archives"].(bool)

	captureDescription := eventhubs.CaptureDescription{
		Enabled: utils.Bool(enabled),
		Encoding: func() *eventhubs.EncodingCaptureDescription {
			v := eventhubs.EncodingCaptureDescription(encoding)
			return &v
		}(),
		IntervalInSeconds: utils.Int64(int64(intervalInSeconds)),
		SizeLimitInBytes:  utils.Int64(int64(sizeLimitInBytes)),
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

			captureDescription.Destination = &eventhubs.Destination{
				Name: utils.String(destinationName),
				Properties: &eventhubs.DestinationProperties{
					ArchiveNameFormat:        utils.String(archiveNameFormat),
					BlobContainer:            utils.String(blobContainerName),
					StorageAccountResourceId: utils.String(storageAccountId),
				},
			}

			if destinationAuthType := destination["storage_authentication"]; destinationAuthType != nil && destinationAuthType.(string) != "" {
				authType := eventhubs.CaptureIdentityType(destinationAuthType.(string))
				captureDescription.Destination.Identity = &eventhubs.CaptureIdentity{
					Type: pointer.To(authType),
				}
			}

			if destinationAuthTypeId := destination["storage_authentication_id"]; destinationAuthTypeId != nil && destinationAuthTypeId.(string) != "" {
				authId := destinationAuthTypeId.(string)
				captureDescription.Destination.Identity.UserAssignedIdentity = &authId
			}
		}
	}

	return &captureDescription
}

func flattenEventHubRetentionDescription(description *eventhubs.RetentionDescription) []interface{} {
	results := make([]interface{}, 0)

	if description != nil {
		output := make(map[string]interface{})

		if cleanupPolicy := description.CleanupPolicy; cleanupPolicy != nil {
			output["cleanup_policy"] = pointer.FromEnum(cleanupPolicy)
		}

		if retentionTimeInHours := description.RetentionTimeInHours; retentionTimeInHours != nil && output["cleanup_policy"] == "Delete" {
			output["retention_time_in_hours"] = pointer.From(retentionTimeInHours)
		}

		if tombstoneRetentionTimeInHours := description.TombstoneRetentionTimeInHours; tombstoneRetentionTimeInHours != nil && output["cleanup_policy"] == "Compact" {
			output["tombstone_retention_time_in_hours"] = pointer.From(tombstoneRetentionTimeInHours)
		}

		results = append(results, output)
	}

	return results
}

func flattenEventHubCaptureDescription(description *eventhubs.CaptureDescription) []interface{} {
	results := make([]interface{}, 0)

	if description != nil {
		output := make(map[string]interface{})

		if enabled := description.Enabled; enabled != nil {
			output["enabled"] = *enabled
		}

		if skipEmptyArchives := description.SkipEmptyArchives; skipEmptyArchives != nil {
			output["skip_empty_archives"] = *skipEmptyArchives
		}

		encoding := ""
		if description.Encoding != nil {
			encoding = string(*description.Encoding)
		}
		output["encoding"] = encoding

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

			if props := destination.Properties; props != nil {
				if archiveNameFormat := props.ArchiveNameFormat; archiveNameFormat != nil {
					destinationOutput["archive_name_format"] = *archiveNameFormat
				}
				if blobContainerName := props.BlobContainer; blobContainerName != nil {
					destinationOutput["blob_container_name"] = *blobContainerName
				}
				if storageAccountId := props.StorageAccountResourceId; storageAccountId != nil {
					destinationOutput["storage_account_id"] = *storageAccountId
				}
			}

			if storageIdentity := destination.Identity; storageIdentity != nil {
				if storageAuthType := storageIdentity.Type; storageAuthType != nil {
					authType := string(pointer.From(storageAuthType))
					destinationOutput["storage_authentication"] = authType
				}
				if storageAuthId := storageIdentity.UserAssignedIdentity; storageAuthId != nil {
					destinationOutput["storage_authentication_id"] = *storageAuthId
				}
			}

			output["destination"] = []interface{}{destinationOutput}
		}

		results = append(results, output)
	}

	return results
}
