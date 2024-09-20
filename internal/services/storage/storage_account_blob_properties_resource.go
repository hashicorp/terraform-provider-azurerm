// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var storageAccountBlobPropertiesResourceName = "azurerm_storage_account_blob_properties"

func resourceStorageAccountBlobProperties() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageAccountBlobPropertiesCreate,
		Read:   resourceStorageAccountBlobPropertiesRead,
		Update: resourceStorageAccountBlobPropertiesUpdate,
		Delete: resourceStorageAccountBlobPropertiesDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"change_feed_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"change_feed_retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 146000),
			},

			"container_delete_retention_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"days": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      7,
							ValidateFunc: validation.IntBetween(1, 365),
						},
					},
				},
			},

			"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

			"default_service_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.BlobPropertiesDefaultServiceVersion,
			},

			"delete_retention_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"days": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      7,
							ValidateFunc: validation.IntBetween(1, 365),
						},

						"permanent_delete_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"last_access_time_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"restore_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"days": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 365),
						},
					},
				},
				RequiredWith: []string{"delete_retention_policy"},
			},

			"versioning_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	return resource
}

func resourceStorageAccountBlobPropertiesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:CREATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	// TODO: Add Import error supported for this resource...

	accountKind := pointer.From(model.Kind)
	dnsEndpointType := pointer.From(model.Properties.DnsEndpointType)
	accountReplicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	accountTier := pointer.From(model.Sku.Tier)
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if !supportLevel.supportBlob {
		return fmt.Errorf("%q are not supported for account kind %q", storageAccountBlobPropertiesResourceName, accountKind)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	accountDetails, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// NOTE: Wait for the blob data plane container to become available...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'custompollers.NewDataPlaneBlobContainersAvailabilityPoller': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	pollerType, err := custompollers.NewDataPlaneBlobContainersAvailabilityPoller(ctx, storageClient, accountDetails)
	if err != nil {
		return fmt.Errorf("building Blob Service Poller: %+v", err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'pollers.NewPoller': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the Blob Service to become available: %+v", err)
	}

	// NOTE: Now that we know the data plane container is available, we can now set the properties on the resource
	// after a bit more validation of the resource...
	props := d.Get("properties").([]interface{})
	blobProperties, err := expandAccountBlobServiceProperties(accountKind, props)
	if err != nil {
		return err
	}

	// See: https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#:~:text=Storage%20accounts%20with%20a%20hierarchical%20namespace%20enabled%20for%20use%20with%20Azure%20Data%20Lake%20Storage%20Gen2%20are%20not%20currently%20supported.
	isVersioningEnabled := pointer.From(blobProperties.Properties.IsVersioningEnabled)
	isHnsEnabled := pointer.From(model.Properties.IsHnsEnabled)
	if isVersioningEnabled && isHnsEnabled {
		return fmt.Errorf("`versioning_enabled` cannot be 'true' when `is_hns_enabled` is true")
	}

	if !isVersioningEnabled {
		if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
			// Otherwise, API returns: "Conflicting feature 'restorePolicy' is enabled. Please disable it and retry."
			return fmt.Errorf("`properties.restore_policy` cannot be set when `versioning_enabled` is 'false'")
		}

		immutableStorageWithVersioningEnabled := false
		if versioning := model.Properties.ImmutableStorageWithVersioning; versioning != nil {
			if versioning.ImmutabilityPolicy != nil && versioning.Enabled != nil {
				immutableStorageWithVersioningEnabled = *versioning.Enabled
			}
		}

		if immutableStorageWithVersioningEnabled {
			// Otherwise, API returns: "Conflicting feature 'Account level WORM' is enabled. Please disable it and retry."
			// See: https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-policy-configure-version-scope?tabs=azure-portal#prerequisites
			return fmt.Errorf("`immutability_policy` cannot be set when `versioning_enabled` is 'false'")
		}
	}

	// TODO: This is a temporary limitation on Storage service. Remove this check once the API supports this scenario.
	// See https://github.com/hashicorp/terraform-provider-azurerm/pull/25450#discussion_r1542471667 for the context.
	if dnsEndpointType == storageaccounts.DnsEndpointTypeAzureDnsZone {
		if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
			// Otherwise, API returns: "Required feature Global Dns is disabled"
			// This is confirmed with the SRP team, where they said:
			// > restorePolicy feature is incompatible with partitioned DNS
			return fmt.Errorf("`properties.restore_policy` cannot be set when `dns_endpoint_type` is set to %q", storageaccounts.DnsEndpointTypeAzureDnsZone)
		}
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.ResourceManager.BlobService.SetServiceProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	if _, err = storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
		return fmt.Errorf("creating `properties`: %+v", err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountBlobPropertiesRead(d, meta)
}

func resourceStorageAccountBlobPropertiesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:UPDATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	accountReplicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	storageType := model.Sku.Name
	dnsEndpointType := pointer.From(model.Properties.DnsEndpointType)
	isHnsEnabled := pointer.From(model.Properties.IsHnsEnabled)
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if accountKind == storageaccounts.KindBlobStorage || accountKind == storageaccounts.KindStorage {
		if storageType == storageaccounts.SkuNameStandardZRS {
			return fmt.Errorf("an `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
		}
	}

	if d.HasChange("properties") {
		if !supportLevel.supportBlob {
			return fmt.Errorf("%q are not supported for account kind %q in sku tier %q", storageAccountBlobPropertiesResourceName, accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, d.Get("properties").([]interface{}))
		if err != nil {
			return err
		}

		if blobProperties.Properties.IsVersioningEnabled != nil && *blobProperties.Properties.IsVersioningEnabled && isHnsEnabled {
			return fmt.Errorf("`versioning_enabled` cannot be true when `is_hns_enabled` is true")
		}

		// Disable restore_policy first. Disabling restore_policy and while setting delete_retention_policy.allow_permanent_delete to true cause error.
		// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
		if v := d.Get("properties.0.restore_policy"); d.HasChange("properties.0.restore_policy") && len(v.([]interface{})) == 0 {
			log.Printf("[DEBUG] [%s:UPDATE] Disabling 'RestorePolicy' prior to changing 'DeleteRetentionPolicy': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
			blobPayload := blobservice.BlobServiceProperties{
				Properties: &blobservice.BlobServicePropertiesProperties{
					RestorePolicy: expandAccountBlobPropertiesRestorePolicy(v.([]interface{})),
				},
			}

			log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.ResourceManager.BlobService.SetServiceProperties' to disable 'RestorePolicy': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
			if _, err := storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, blobPayload); err != nil {
				return fmt.Errorf("updating Azure Storage Account blob restore policy %q: %+v", id.StorageAccountName, err)
			}
		}

		if dnsEndpointType == storageaccounts.DnsEndpointTypeAzureDnsZone {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`properties.restore_policy` cannot be set when `dns_endpoint_type` is set to %q", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.ResourceManager.BlobService.SetServiceProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
		if _, err = storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
			return fmt.Errorf("updating `properties` for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountBlobPropertiesRead(d, meta)
}

func resourceStorageAccountBlobPropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	// we then need to find the storage account
	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.ResourceManager.BlobService.GetServiceProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	blobProps, err := storageClient.ResourceManager.BlobService.GetServiceProperties(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading blob properties for %s: %+v", *id, err)
	}

	blobProperties := flattenAccountBlobServiceProperties(blobProps.Model)

	if err := d.Set("properties", blobProperties); err != nil {
		return fmt.Errorf("setting `properties` for %s: %+v", *id, err)
	}

	return nil
}

func resourceStorageAccountBlobPropertiesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:DELETE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountKind := pointer.From(model.Kind)
	blobProperties, err := expandAccountBlobServiceProperties(accountKind, make([]interface{}, 0))
	if err != nil {
		return err
	}

	// NOTE: You need to disable the 'restore_policy' before you can set the 'delete_retention_policy.allow_permanent_delete' to 'true'
	// else you will get the "ConflictFeatureEnabled: Conflicting feature 'restorePolicy' is 'enabled' returned by the RP.
	// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
	blobPayload := blobservice.BlobServiceProperties{
		Properties: &blobservice.BlobServicePropertiesProperties{
			RestorePolicy: expandAccountBlobPropertiesRestorePolicy(make([]interface{}, 0)),
		},
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.ResourceManager.BlobService.SetServiceProperties' to disable 'RestorePolicy': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	if _, err := storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, blobPayload); err != nil {
		return fmt.Errorf("updating the Azure Storage Account %q 'RestorePolicy' prior to the deletion of %q: %+v", id.StorageAccountName, storageAccountBlobPropertiesResourceName, err)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.ResourceManager.BlobService.SetServiceProperties': %s", strings.ToUpper(storageAccountBlobPropertiesResourceName), id)
	if _, err = storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
		return fmt.Errorf("deleting %q for %s: %+v", storageAccountBlobPropertiesResourceName, *id, err)
	}

	return nil
}

func expandAccountBlobServiceProperties(kind storageaccounts.Kind, input []interface{}) (*blobservice.BlobServiceProperties, error) {
	props := blobservice.BlobServicePropertiesProperties{
		Cors: &blobservice.CorsRules{
			CorsRules: &[]blobservice.CorsRule{},
		},
		DeleteRetentionPolicy: &blobservice.DeleteRetentionPolicy{
			Enabled: utils.Bool(false),
		},
	}

	// `Storage` (v1) kind doesn't support:
	// - LastAccessTimeTrackingPolicy: Confirmed by SRP.
	// - ChangeFeed: See https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blob-change-feed?tabs=azure-portal#enable-and-disable-the-change-feed.
	// - Versioning: See https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#how-blob-versioning-works
	// - Restore Policy: See https://learn.microsoft.com/en-us/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
	if kind != storageaccounts.KindStorage {
		props.LastAccessTimeTrackingPolicy = &blobservice.LastAccessTimeTrackingPolicy{
			Enable: false,
		}
		props.ChangeFeed = &blobservice.ChangeFeed{
			Enabled: pointer.To(false),
		}
		props.IsVersioningEnabled = pointer.To(false)
	}

	if len(input) > 0 {
		v := input[0].(map[string]interface{})

		deletePolicyRaw := v["delete_retention_policy"].([]interface{})
		props.DeleteRetentionPolicy = expandAccountBlobDeleteRetentionPolicy(deletePolicyRaw)

		containerDeletePolicyRaw := v["container_delete_retention_policy"].([]interface{})
		props.ContainerDeleteRetentionPolicy = expandAccountBlobContainerDeleteRetentionPolicy(containerDeletePolicyRaw)

		corsRaw := v["cors_rule"].([]interface{})
		props.Cors = expandAccountBlobPropertiesCors(corsRaw)

		props.IsVersioningEnabled = pointer.To(v["versioning_enabled"].(bool))

		if version, ok := v["default_service_version"].(string); ok && version != "" {
			props.DefaultServiceVersion = pointer.To(version)
		}

		// `Storage` (v1) kind doesn't support:
		// - LastAccessTimeTrackingPolicy
		// - ChangeFeed
		// - Versioning
		// - RestorePolicy
		lastAccessTimeEnabled := v["last_access_time_enabled"].(bool)
		changeFeedEnabled := v["change_feed_enabled"].(bool)
		changeFeedRetentionInDays := v["change_feed_retention_in_days"].(int)
		restorePolicyRaw := v["restore_policy"].([]interface{})
		versioningEnabled := v["versioning_enabled"].(bool)
		if kind != storageaccounts.KindStorage {
			props.LastAccessTimeTrackingPolicy = &blobservice.LastAccessTimeTrackingPolicy{
				Enable: lastAccessTimeEnabled,
			}
			props.ChangeFeed = &blobservice.ChangeFeed{
				Enabled: pointer.To(changeFeedEnabled),
			}
			if changeFeedRetentionInDays != 0 {
				props.ChangeFeed.RetentionInDays = pointer.To(int64(changeFeedRetentionInDays))
			}
			props.RestorePolicy = expandAccountBlobPropertiesRestorePolicy(restorePolicyRaw)
			props.IsVersioningEnabled = &versioningEnabled
		} else {
			if lastAccessTimeEnabled {
				return nil, fmt.Errorf("`last_access_time_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedEnabled {
				return nil, fmt.Errorf("`change_feed_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedRetentionInDays != 0 {
				return nil, fmt.Errorf("`change_feed_retention_in_days` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if len(restorePolicyRaw) != 0 {
				return nil, fmt.Errorf("`restore_policy` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if versioningEnabled {
				return nil, fmt.Errorf("`versioning_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
		}

		// Sanity check for the prerequisites of restore_policy
		// Ref: https://learn.microsoft.com/en-us/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
		if p := props.RestorePolicy; p != nil && p.Enabled {
			if props.ChangeFeed == nil || props.ChangeFeed.Enabled == nil || !*props.ChangeFeed.Enabled {
				return nil, fmt.Errorf("`change_feed_enabled` must be `true` when `restore_policy` is set")
			}
			if props.IsVersioningEnabled == nil || !*props.IsVersioningEnabled {
				return nil, fmt.Errorf("`versioning_enabled` must be `true` when `restore_policy` is set")
			}
		}
	}

	return &blobservice.BlobServiceProperties{
		Properties: &props,
	}, nil
}

func flattenAccountBlobServiceProperties(input *blobservice.BlobServiceProperties) []interface{} {
	if input == nil || input.Properties == nil {
		return []interface{}{}
	}

	flattenedCorsRules := make([]interface{}, 0)
	if corsRules := input.Properties.Cors; corsRules != nil {
		flattenedCorsRules = flattenAccountBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]interface{}, 0)
	if deletePolicy := input.Properties.DeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenAccountBlobDeleteRetentionPolicy(deletePolicy)
	}

	flattenedRestorePolicy := make([]interface{}, 0)
	if restorePolicy := input.Properties.RestorePolicy; restorePolicy != nil {
		flattenedRestorePolicy = flattenAccountBlobPropertiesRestorePolicy(restorePolicy)
	}

	flattenedContainerDeletePolicy := make([]interface{}, 0)
	if containerDeletePolicy := input.Properties.ContainerDeleteRetentionPolicy; containerDeletePolicy != nil {
		flattenedContainerDeletePolicy = flattenAccountBlobContainerDeleteRetentionPolicy(containerDeletePolicy)
	}

	versioning, changeFeedEnabled, changeFeedRetentionInDays := false, false, 0
	if input.Properties.IsVersioningEnabled != nil {
		versioning = *input.Properties.IsVersioningEnabled
	}

	if v := input.Properties.ChangeFeed; v != nil {
		if v.Enabled != nil {
			changeFeedEnabled = *v.Enabled
		}
		if v.RetentionInDays != nil {
			changeFeedRetentionInDays = int(*v.RetentionInDays)
		}
	}

	var defaultServiceVersion string
	if input.Properties.DefaultServiceVersion != nil {
		defaultServiceVersion = *input.Properties.DefaultServiceVersion
	}

	var LastAccessTimeTrackingPolicy bool
	if v := input.Properties.LastAccessTimeTrackingPolicy; v != nil {
		LastAccessTimeTrackingPolicy = v.Enable
	}

	return []interface{}{
		map[string]interface{}{
			"change_feed_enabled":               changeFeedEnabled,
			"change_feed_retention_in_days":     changeFeedRetentionInDays,
			"container_delete_retention_policy": flattenedContainerDeletePolicy,
			"cors_rule":                         flattenedCorsRules,
			"default_service_version":           defaultServiceVersion,
			"delete_retention_policy":           flattenedDeletePolicy,
			"last_access_time_enabled":          LastAccessTimeTrackingPolicy,
			"restore_policy":                    flattenedRestorePolicy,
			"versioning_enabled":                versioning,
		},
	}
}

func expandAccountBlobDeleteRetentionPolicy(input []interface{}) *blobservice.DeleteRetentionPolicy {
	result := blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.DeleteRetentionPolicy{
		Enabled:              pointer.To(true),
		AllowPermanentDelete: pointer.To(policy["permanent_delete_enabled"].(bool)),
		Days:                 pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobDeleteRetentionPolicy(input *blobservice.DeleteRetentionPolicy) []interface{} {
	deleteRetentionPolicy := make([]interface{}, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		var permanentDeleteEnabled bool
		if input.AllowPermanentDelete != nil {
			permanentDeleteEnabled = *input.AllowPermanentDelete
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]interface{}{
			"days":                     days,
			"permanent_delete_enabled": permanentDeleteEnabled,
		})
	}

	return deleteRetentionPolicy
}

func expandAccountBlobContainerDeleteRetentionPolicy(input []interface{}) *blobservice.DeleteRetentionPolicy {
	result := blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobContainerDeleteRetentionPolicy(input *blobservice.DeleteRetentionPolicy) []interface{} {
	deleteRetentionPolicy := make([]interface{}, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]interface{}{
			"days": days,
		})
	}

	return deleteRetentionPolicy
}

func expandAccountBlobPropertiesRestorePolicy(input []interface{}) *blobservice.RestorePolicyProperties {
	result := blobservice.RestorePolicyProperties{
		Enabled: false,
	}

	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.RestorePolicyProperties{
		Enabled: true,
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobPropertiesRestorePolicy(input *blobservice.RestorePolicyProperties) []interface{} {
	restorePolicy := make([]interface{}, 0)

	if input == nil {
		return restorePolicy
	}

	if enabled := input.Enabled; enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		restorePolicy = append(restorePolicy, map[string]interface{}{
			"days": days,
		})
	}

	return restorePolicy
}

func expandAccountBlobPropertiesCors(input []interface{}) *blobservice.CorsRules {
	blobCorsRules := blobservice.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]blobservice.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]interface{})

			allowedMethods := make([]blobservice.AllowedMethods, 0)
			for _, val := range *utils.ExpandStringSlice(item["allowed_methods"].([]interface{})) {
				allowedMethods = append(allowedMethods, blobservice.AllowedMethods(val))
			}
			corsRules = append(corsRules, blobservice.CorsRule{
				AllowedHeaders:  *utils.ExpandStringSlice(item["allowed_headers"].([]interface{})),
				AllowedOrigins:  *utils.ExpandStringSlice(item["allowed_origins"].([]interface{})),
				AllowedMethods:  allowedMethods,
				ExposedHeaders:  *utils.ExpandStringSlice(item["exposed_headers"].([]interface{})),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func flattenAccountBlobPropertiesCorsRule(input *blobservice.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}
