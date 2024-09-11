// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var storageAccountSharePropertiesResourceName = "azurerm_storage_account_share_properties"

func resourceStorageAccountShareProperties() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageAccountSharePropertiesCreate,
		Read:   resourceStorageAccountSharePropertiesRead,
		Update: resourceStorageAccountSharePropertiesUpdate,
		Delete: resourceStorageAccountSharePropertiesDelete,

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

			"properties": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

						"retention_policy": {
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

						"smb": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"authentication_types": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Kerberos",
												"NTLMv2",
											}, false),
										},
									},

									"channel_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"AES-128-CCM",
												"AES-128-GCM",
												"AES-256-GCM",
											}, false),
										},
									},

									"kerberos_ticket_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"AES-256",
												"RC4-HMAC",
											}, false),
										},
									},

									"multichannel_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"versions": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"SMB2.1",
												"SMB3.0",
												"SMB3.1.1",
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return resource
}

func resourceStorageAccountSharePropertiesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	// TODO: add import check here...
	// if !response.WasNotFound(existing.HttpResponse) {
	// 	return tf.ImportAsExistsError(storageAccountResourceName, id.ID())
	// }

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	fileShareSupportedReplicationTypes := []string{"LRS", "GRS", "RAGRS"}

	// File share is only supported for StorageV2 and FileStorage, StorageV2 does not support file share if the StorageV2 is a 'Premium' Sku Tier.
	// See: https://docs.microsoft.com/en-us/azure/storage/files/storage-files-planning#management-concepts
	supportsShare := (accountKind == storageaccounts.KindFileStorage ||
		(accountTier != storageaccounts.SkuTierPremium && (accountKind == storageaccounts.KindStorageVTwo ||
			(accountKind == storageaccounts.KindStorage && slices.Contains(fileShareSupportedReplicationTypes, replicationType)))))
	// Per local testing, Storage V1 accounts (e.g., storageaccounts.KindStorage) with a replication type of LRS/GRS/RAGRS do support file endpoints,
	// while GZRS and RAGZRS are invalid and ZRS is valid but it does not have a file endpoint...

	if !supportsShare {
		return fmt.Errorf("%q are not supported for account kind %q in sku tier %q", storageAccountSharePropertiesResourceName, accountKind, accountTier)
	}

	// NOTE: Wait for the data plane share container to become available...
	log.Printf("[DEBUG] [CREATE] Calling 'storageClient.FindAccount' for %s", id)
	accountDetails, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if accountDetails == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	log.Printf("[DEBUG] [CREATE] Calling 'custompollers.NewDataPlaneFileShareAvailabilityPoller' building File Share Poller for %s", id)
	pollerType, err := custompollers.NewDataPlaneFileShareAvailabilityPoller(storageClient, accountDetails)
	if err != nil {
		return fmt.Errorf("building File Share Poller: %+v", err)
	}

	log.Printf("[DEBUG] [CREATE] Calling 'poller.PollUntilDone' waiting for the File Service to become available for %s", id)
	poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the File Service to become available: %+v", err)
	}

	// NOTE: Now that we know the data plane container is available, we can now set the properties on the resource
	// after a bit more validation of the resource...
	sharePayload := expandAccountShareProperties(d.Get("properties").([]interface{}))

	// The API complains if any multichannel info is sent on non premium file shares. Even if multichannel is set to false
	if accountTier != storageaccounts.SkuTierPremium && sharePayload.Properties != nil && sharePayload.Properties.ProtocolSettings != nil {
		// Error if the user has tried to enable multichannel on a standard tier storage account
		smb := sharePayload.Properties.ProtocolSettings.Smb
		if smb != nil && smb.Multichannel != nil {
			if smb.Multichannel.Enabled != nil && *smb.Multichannel.Enabled {
				return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
			}

			sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
		}
	}

	log.Printf("[DEBUG] [CREATE] Calling 'storageClient.ResourceManager.FileService.SetServiceProperties' for %s", id)
	if _, err = storageClient.ResourceManager.FileService.SetServiceProperties(ctx, *id, sharePayload); err != nil {
		return fmt.Errorf("creating %q: %+v", storageAccountSharePropertiesResourceName, err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountSharePropertiesRead(d, meta)
}

func resourceStorageAccountSharePropertiesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	keyVaultClient := meta.(*clients.Client).KeyVault
	dataPlaneEnabled := meta.(*clients.Client).Features.Storage.DataPlaneAccessOnReadEnabled
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	accountTier := storageaccounts.SkuTier(d.Get("account_tier").(string))
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := storageaccounts.Kind(d.Get("account_kind").(string))

	if accountKind == storageaccounts.KindBlobStorage || accountKind == storageaccounts.KindStorage {
		if storageType == string(storageaccounts.SkuNameStandardZRS) {
			return fmt.Errorf("an `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
		}
	}

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}
	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	props := storageaccounts.StorageAccountPropertiesCreateParameters{
		AccessTier:                            existing.Model.Properties.AccessTier,
		AllowBlobPublicAccess:                 existing.Model.Properties.AllowBlobPublicAccess,
		AllowedCopyScope:                      existing.Model.Properties.AllowedCopyScope,
		AllowSharedKeyAccess:                  existing.Model.Properties.AllowSharedKeyAccess,
		AllowCrossTenantReplication:           existing.Model.Properties.AllowCrossTenantReplication,
		AzureFilesIdentityBasedAuthentication: existing.Model.Properties.AzureFilesIdentityBasedAuthentication,
		CustomDomain:                          existing.Model.Properties.CustomDomain,
		DefaultToOAuthAuthentication:          existing.Model.Properties.DefaultToOAuthAuthentication,
		DnsEndpointType:                       existing.Model.Properties.DnsEndpointType,
		Encryption:                            existing.Model.Properties.Encryption,
		KeyPolicy:                             existing.Model.Properties.KeyPolicy,
		ImmutableStorageWithVersioning:        existing.Model.Properties.ImmutableStorageWithVersioning,
		IsNfsV3Enabled:                        existing.Model.Properties.IsNfsV3Enabled,
		IsSftpEnabled:                         existing.Model.Properties.IsSftpEnabled,
		IsLocalUserEnabled:                    existing.Model.Properties.IsLocalUserEnabled,
		IsHnsEnabled:                          existing.Model.Properties.IsHnsEnabled,
		MinimumTlsVersion:                     existing.Model.Properties.MinimumTlsVersion,
		NetworkAcls:                           existing.Model.Properties.NetworkAcls,
		PublicNetworkAccess:                   existing.Model.Properties.PublicNetworkAccess,
		RoutingPreference:                     existing.Model.Properties.RoutingPreference,
		SasPolicy:                             existing.Model.Properties.SasPolicy,
		SupportsHTTPSTrafficOnly:              existing.Model.Properties.SupportsHTTPSTrafficOnly,
	}

	if existing.Model.Properties.LargeFileSharesState != nil && *existing.Model.Properties.LargeFileSharesState == storageaccounts.LargeFileSharesStateEnabled {
		// We can only set this if it's Enabled, else the API complains during Update that we're sending Disabled, even if it's always been off
		props.LargeFileSharesState = existing.Model.Properties.LargeFileSharesState
	}

	expandedIdentity := existing.Model.Identity
	if d.HasChange("identity") {
		expandedIdentity, err = identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
	}

	if d.HasChange("access_tier") {
		props.AccessTier = pointer.To(storageaccounts.AccessTier(d.Get("access_tier").(string)))
	}
	if d.HasChange("allowed_copy_scope") {
		props.AllowedCopyScope = pointer.To(storageaccounts.AllowedCopyScope(d.Get("allowed_copy_scope").(string)))
	}
	if d.HasChange("allow_nested_items_to_be_public") {
		props.AllowBlobPublicAccess = pointer.To(d.Get("allow_nested_items_to_be_public").(bool))
	}
	if d.HasChange("cross_tenant_replication_enabled") {
		props.AllowCrossTenantReplication = pointer.To(d.Get("cross_tenant_replication_enabled").(bool))
	}
	if d.HasChange("custom_domain") {
		props.CustomDomain = expandAccountCustomDomain(d.Get("custom_domain").([]interface{}))
	}
	if d.HasChange("customer_managed_key") {
		queueEncryptionKeyType := storageaccounts.KeyType(d.Get("queue_encryption_key_type").(string))
		tableEncryptionKeyType := storageaccounts.KeyType(d.Get("table_encryption_key_type").(string))
		encryptionRaw := d.Get("customer_managed_key").([]interface{})
		encryption, err := expandAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, encryptionRaw, accountTier, accountKind, *expandedIdentity, queueEncryptionKeyType, tableEncryptionKeyType)
		if err != nil {
			return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
		}
		props.Encryption = encryption
	}
	if d.HasChange("shared_access_key_enabled") {
		props.AllowSharedKeyAccess = pointer.To(d.Get("shared_access_key_enabled").(bool))
	} else {
		// If AllowSharedKeyAccess is nil that breaks the Portal UI as reported in https://github.com/hashicorp/terraform-provider-azurerm/issues/11689
		// currently the Portal UI reports nil as false, and per the ARM API documentation nil is true. This manifests itself in the Portal UI
		// when a storage account is created by terraform that the AllowSharedKeyAccess is Disabled when it is actually Enabled, thus confusing out customers
		// to fix this, I have added this code to explicitly to set the value to true if is nil to workaround the Portal UI bug for our customers.
		// this is designed as a passive change, meaning the change will only take effect when the existing storage account is modified in some way if the
		// account already exists. since I have also switched up the default behaviour for net new storage accounts to always set this value as true, this issue
		// should automatically correct itself over time with these changes.
		// TODO: Remove code when Portal UI team fixes their code
		if sharedKeyAccess := props.AllowSharedKeyAccess; sharedKeyAccess == nil {
			props.AllowSharedKeyAccess = pointer.To(true)
		}
	}
	if d.HasChange("default_to_oauth_authentication") {
		props.DefaultToOAuthAuthentication = pointer.To(d.Get("default_to_oauth_authentication").(bool))
	}

	if d.HasChange("https_traffic_only_enabled") {
		props.SupportsHTTPSTrafficOnly = pointer.To(d.Get("https_traffic_only_enabled").(bool))
	}

	if !features.FourPointOhBeta() {
		if d.HasChange("enable_https_traffic_only") {
			props.SupportsHTTPSTrafficOnly = pointer.To(d.Get("enable_https_traffic_only").(bool))
		}
	}

	if d.HasChange("large_file_share_enabled") {
		// largeFileSharesState can only be set to `Enabled` and not `Disabled`, even if it is currently `Disabled`
		if oldValue, newValue := d.GetChange("large_file_share_enabled"); oldValue.(bool) && !newValue.(bool) {
			return fmt.Errorf("`large_file_share_enabled` cannot be disabled once it's been enabled")
		}

		if _, ok := storageKindsSupportLargeFileShares[accountKind]; !ok {
			keys := sortedKeysFromSlice(storageKindsSupportLargeFileShares)
			return fmt.Errorf("`large_file_shares_enabled` can only be set to `true` with `account_kind` set to one of: %+v", strings.Join(keys, " / "))
		}
		props.LargeFileSharesState = pointer.To(storageaccounts.LargeFileSharesStateEnabled)
	}

	if d.HasChange("local_user_enabled") {
		props.IsLocalUserEnabled = pointer.To(d.Get("local_user_enabled").(bool))
	}

	if d.HasChange("min_tls_version") {
		props.MinimumTlsVersion = pointer.To(storageaccounts.MinimumTlsVersion(d.Get("min_tls_version").(string)))
	}

	if d.HasChange("network_rules") {
		props.NetworkAcls = expandAccountNetworkRules(d.Get("network_rules").([]interface{}), tenantId)
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := storageaccounts.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = storageaccounts.PublicNetworkAccessEnabled
		}
		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("routing") {
		props.RoutingPreference = expandAccountRoutingPreference(d.Get("routing").([]interface{}))
	}

	if d.HasChange("sas_policy") {
		// TODO: Currently, there is no way to represent a `null` value in the payload - instead it will be omitted, `sas_policy` can not be disabled once enabled.
		props.SasPolicy = expandAccountSASPolicy(d.Get("sas_policy").([]interface{}))
	}

	if d.HasChange("sftp_enabled") {
		props.IsSftpEnabled = pointer.To(d.Get("sftp_enabled").(bool))
	}

	payload := storageaccounts.StorageAccountCreateParameters{
		ExtendedLocation: existing.Model.ExtendedLocation,
		Kind:             *existing.Model.Kind,
		Location:         existing.Model.Location,
		Identity:         existing.Model.Identity,
		Properties:       &props,
		Sku:              *existing.Model.Sku,
		Tags:             existing.Model.Tags,
	}

	// ensure any top-level properties are updated
	if d.HasChange("account_kind") {
		payload.Kind = accountKind
	}

	if d.HasChange("account_replication_type") {
		// storageType is derived from "account_replication_type" and "account_tier" (force-new)
		payload.Sku = storageaccounts.Sku{
			Name: storageaccounts.SkuName(storageType),
		}
	}

	if d.HasChange("identity") {
		payload.Identity = expandedIdentity
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// azure_files_authentication must be the last to be updated, cause it'll occupy the storage account for several minutes after receiving the response 200 OK. Issue: https://github.com/Azure/azure-rest-api-specs/issues/11272
	if d.HasChange("azure_files_authentication") {
		// due to service issue: https://github.com/Azure/azure-rest-api-specs/issues/12473, we need to update to None before changing its DirectoryServiceOptions
		old, new := d.GetChange("azure_files_authentication.0.directory_type")
		if old != new && new != string(storageaccounts.DirectoryServiceOptionsNone) {
			log.Print("[DEBUG] Disabling AzureFilesIdentityBasedAuthentication prior to changing DirectoryServiceOptions")
			dsNone := storageaccounts.StorageAccountUpdateParameters{
				Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
					AzureFilesIdentityBasedAuthentication: &storageaccounts.AzureFilesIdentityBasedAuthentication{
						DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsNone,
					},
				},
			}

			if _, err := client.Update(ctx, *id, dsNone); err != nil {
				return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
			}
		}

		expandAADFilesAuthentication, err := expandAccountAzureFilesAuthentication(d.Get("azure_files_authentication").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `azure_files_authentication`: %+v", err)
		}
		opts := storageaccounts.StorageAccountUpdateParameters{
			Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
				AzureFilesIdentityBasedAuthentication: expandAADFilesAuthentication,
			},
		}

		if _, err := client.Update(ctx, *id, opts); err != nil {
			return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
		}
	}

	// Followings are updates to the sub-services
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if d.HasChange("blob_properties") {
		if !supportLevel.supportBlob {
			return fmt.Errorf("`blob_properties` are not supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, d.Get("blob_properties").([]interface{}))
		if err != nil {
			return err
		}

		if blobProperties.Properties.IsVersioningEnabled != nil && *blobProperties.Properties.IsVersioningEnabled && d.Get("is_hns_enabled").(bool) {
			return fmt.Errorf("`versioning_enabled` cannot be true when `is_hns_enabled` is true")
		}

		// Disable restore_policy first. Disabling restore_policy and while setting delete_retention_policy.allow_permanent_delete to true cause error.
		// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
		if v := d.Get("blob_properties.0.restore_policy"); d.HasChange("blob_properties.0.restore_policy") && len(v.([]interface{})) == 0 {
			log.Print("[DEBUG] Disabling RestorePolicy prior to changing DeleteRetentionPolicy")
			blobPayload := blobservice.BlobServiceProperties{
				Properties: &blobservice.BlobServicePropertiesProperties{
					RestorePolicy: expandAccountBlobPropertiesRestorePolicy(v.([]interface{})),
				},
			}

			if _, err := storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, blobPayload); err != nil {
				return fmt.Errorf("updating Azure Storage Account blob restore policy %q: %+v", id.StorageAccountName, err)
			}
		}

		if d.Get("dns_endpoint_type").(string) == string(storageaccounts.DnsEndpointTypeAzureDnsZone) {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` cannot be set when `dns_endpoint_type` is set to `%s`", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		if _, err = storageClient.ResourceManager.BlobService.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties` for %s: %+v", *id, err)
		}
	}

	if !features.FourPointOhBeta() {
		if dataPlaneEnabled {
			if d.HasChange("queue_properties") {
				if !supportLevel.supportQueue {
					return fmt.Errorf("`queue_properties` are not supported for account kind %q in sku tier %q", accountKind, accountTier)
				}

				account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}
				if account == nil {
					return fmt.Errorf("unable to locate %s", *id)
				}

				queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
				if err != nil {
					return fmt.Errorf("building Queues Client: %s", err)
				}

				queueProperties, err := expandAccountQueueProperties(d.Get("queue_properties").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `queue_properties` for %s: %+v", *id, err)
				}

				if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
					return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
				}
			}

			if d.HasChange("static_website") {
				if !supportLevel.supportStaticWebsite {
					return fmt.Errorf("`static_website` are not supported for account kind %q in sku tier %q", accountKind, accountTier)
				}

				account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}
				if account == nil {
					return fmt.Errorf("unable to locate %s", *id)
				}

				accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
				if err != nil {
					return fmt.Errorf("building Data Plane client for %s: %+v", *id, err)
				}

				staticWebsiteProps := expandAccountStaticWebsiteProperties(d.Get("static_website").([]interface{}))

				if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
					return fmt.Errorf("updating `static_website` for %s: %+v", *id, err)
				}
			}
		} else {
			log.Printf("[DEBUG] storage account update for 'queue_properties' and 'static_website' skipped due to DataPlaneAccessOnReadEnabled feature flag being set to 'false'.")
		}
	}

	if d.HasChange("share_properties") {
		if !supportLevel.supportShare {
			return fmt.Errorf("`share_properties` are not supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		sharePayload := expandAccountShareProperties(d.Get("share_properties").([]interface{}))
		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			if sharePayload.Properties.ProtocolSettings.Smb.Multichannel != nil && sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled != nil {
				if *sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled {
					return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}
			}

			sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
		}

		if _, err = storageClient.ResourceManager.FileService.SetServiceProperties(ctx, *id, sharePayload); err != nil {
			return fmt.Errorf("updating File Share Properties for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountSharePropertiesRead(d, meta)
}

func resourceStorageAccountSharePropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	env := meta.(*clients.Client).Account.Environment
	dataPlaneEnabled := meta.(*clients.Client).Features.Storage.DataPlaneAccessOnReadEnabled
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// we then need to find the storage account
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	listKeysOpts := storageaccounts.DefaultListKeysOperationOptions()
	listKeysOpts.Expand = pointer.To(storageaccounts.ListKeyExpandKerb)
	keys, err := client.ListKeys(ctx, *id, listKeysOpts)
	if err != nil {
		hasWriteLock := response.WasConflict(keys.HttpResponse)
		doesntHavePermissions := response.WasForbidden(keys.HttpResponse) || response.WasStatusCode(keys.HttpResponse, http.StatusUnauthorized)
		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	d.Set("name", id.StorageAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	supportLevel := storageAccountServiceSupportLevel{
		supportBlob:          false,
		supportQueue:         false,
		supportShare:         false,
		supportStaticWebsite: false,
	}
	var accountKind storageaccounts.Kind
	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference
	if model := resp.Model; model != nil {
		if model.Kind != nil {
			accountKind = *model.Kind
		}
		d.Set("account_kind", string(accountKind))

		var accountTier storageaccounts.SkuTier
		accountReplicationType := ""
		if sku := model.Sku; sku != nil {
			accountReplicationType = strings.Split(string(sku.Name), "_")[1]
			if sku.Tier != nil {
				accountTier = *sku.Tier
			}
		}
		d.Set("account_tier", string(accountTier))
		d.Set("account_replication_type", accountReplicationType)

		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			primaryEndpoints = props.PrimaryEndpoints
			routingPreference = props.RoutingPreference
			secondaryEndpoints = props.SecondaryEndpoints

			d.Set("access_tier", pointer.From(props.AccessTier))
			d.Set("allowed_copy_scope", pointer.From(props.AllowedCopyScope))
			if err := d.Set("azure_files_authentication", flattenAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
				return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
			}
			d.Set("cross_tenant_replication_enabled", pointer.From(props.AllowCrossTenantReplication))
			d.Set("https_traffic_only_enabled", pointer.From(props.SupportsHTTPSTrafficOnly))
			if !features.FourPointOhBeta() {
				d.Set("enable_https_traffic_only", pointer.From(props.SupportsHTTPSTrafficOnly))
			}
			d.Set("is_hns_enabled", pointer.From(props.IsHnsEnabled))
			d.Set("nfsv3_enabled", pointer.From(props.IsNfsV3Enabled))
			d.Set("primary_location", pointer.From(props.PrimaryLocation))
			if err := d.Set("routing", flattenAccountRoutingPreference(props.RoutingPreference)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
			d.Set("secondary_location", pointer.From(props.SecondaryLocation))
			d.Set("sftp_enabled", pointer.From(props.IsSftpEnabled))

			// NOTE: The Storage API returns `null` rather than the default value in the API response for existing
			// resources when a new field gets added - meaning we need to default the values below.
			allowBlobPublicAccess := true
			if props.AllowBlobPublicAccess != nil {
				allowBlobPublicAccess = *props.AllowBlobPublicAccess
			}
			d.Set("allow_nested_items_to_be_public", allowBlobPublicAccess)

			defaultToOAuthAuthentication := false
			if props.DefaultToOAuthAuthentication != nil {
				defaultToOAuthAuthentication = *props.DefaultToOAuthAuthentication
			}
			d.Set("default_to_oauth_authentication", defaultToOAuthAuthentication)

			dnsEndpointType := storageaccounts.DnsEndpointTypeStandard
			if props.DnsEndpointType != nil {
				dnsEndpointType = *props.DnsEndpointType
			}
			d.Set("dns_endpoint_type", dnsEndpointType)

			isLocalEnabled := true
			if props.IsLocalUserEnabled != nil {
				isLocalEnabled = *props.IsLocalUserEnabled
			}
			d.Set("local_user_enabled", isLocalEnabled)

			largeFileShareEnabled := false
			if props.LargeFileSharesState != nil {
				largeFileShareEnabled = *props.LargeFileSharesState == storageaccounts.LargeFileSharesStateEnabled
			}
			d.Set("large_file_share_enabled", largeFileShareEnabled)

			minTlsVersion := string(storageaccounts.MinimumTlsVersionTLSOneZero)
			if props.MinimumTlsVersion != nil {
				minTlsVersion = string(*props.MinimumTlsVersion)
			}
			d.Set("min_tls_version", minTlsVersion)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == storageaccounts.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			allowSharedKeyAccess := true
			if props.AllowSharedKeyAccess != nil {
				allowSharedKeyAccess = *props.AllowSharedKeyAccess
			}
			d.Set("shared_access_key_enabled", allowSharedKeyAccess)

			if err := d.Set("custom_domain", flattenAccountCustomDomain(props.CustomDomain)); err != nil {
				return fmt.Errorf("setting `custom_domain`: %+v", err)
			}
			if err := d.Set("immutability_policy", flattenAccountImmutabilityPolicy(props.ImmutableStorageWithVersioning)); err != nil {
				return fmt.Errorf("setting `immutability_policy`: %+v", err)
			}
			if err := d.Set("network_rules", flattenAccountNetworkRules(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_rules`: %+v", err)
			}

			// When the encryption key type is "Service", the queue/table is not returned in the service list, so we default
			// the encryption key type to "Service" if it is absent (must also be the default value for "Service" in the schema)
			infrastructureEncryption := false
			queueEncryptionKeyType := string(storageaccounts.KeyTypeService)
			tableEncryptionKeyType := string(storageaccounts.KeyTypeService)
			if encryption := props.Encryption; encryption != nil {
				infrastructureEncryption = pointer.From(encryption.RequireInfrastructureEncryption)
				if encryption.Services != nil {
					if encryption.Services.Queue != nil && encryption.Services.Queue.KeyType != nil {
						queueEncryptionKeyType = string(*encryption.Services.Queue.KeyType)
					}
					if encryption.Services.Table != nil && encryption.Services.Table.KeyType != nil {
						tableEncryptionKeyType = string(*encryption.Services.Table.KeyType)
					}
				}
			}
			d.Set("infrastructure_encryption_enabled", infrastructureEncryption)
			d.Set("queue_encryption_key_type", queueEncryptionKeyType)
			d.Set("table_encryption_key_type", tableEncryptionKeyType)

			customerManagedKey := flattenAccountCustomerManagedKey(props.Encryption, env)
			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			if err := d.Set("sas_policy", flattenAccountSASPolicy(props.SasPolicy)); err != nil {
				return fmt.Errorf("setting `sas_policy`: %+v", err)
			}

			supportLevel = availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)
		}

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	endpoints := flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	if err := endpoints.set(d); err != nil {
		return err
	}

	storageAccountKeys := make([]storageaccounts.StorageAccountKey, 0)
	if keys.Model != nil && keys.Model.Keys != nil {
		storageAccountKeys = *keys.Model.Keys
	}
	keysAndConnectionStrings := flattenAccountAccessKeysAndConnectionStrings(id.StorageAccountName, *storageDomainSuffix, storageAccountKeys, endpoints)
	if err := keysAndConnectionStrings.set(d); err != nil {
		return err
	}

	blobProperties := make([]interface{}, 0)
	if supportLevel.supportBlob {
		blobProps, err := storageClient.ResourceManager.BlobService.GetServiceProperties(ctx, *id)
		if err != nil {
			return fmt.Errorf("reading blob properties for %s: %+v", *id, err)
		}

		blobProperties = flattenAccountBlobServiceProperties(blobProps.Model)
	}
	if err := d.Set("blob_properties", blobProperties); err != nil {
		return fmt.Errorf("setting `blob_properties` for %s: %+v", *id, err)
	}

	shareProperties := make([]interface{}, 0)
	if supportLevel.supportShare {
		shareProps, err := storageClient.ResourceManager.FileService.GetServiceProperties(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving share properties for %s: %+v", *id, err)
		}

		shareProperties = flattenAccountShareProperties(shareProps.Model)
	}
	if err := d.Set("share_properties", shareProperties); err != nil {
		return fmt.Errorf("setting `share_properties` for %s: %+v", *id, err)
	}

	if dataPlaneEnabled {
		queueProperties := make([]interface{}, 0)
		if supportLevel.supportQueue {
			queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client: %s", err)
			}

			queueProps, err := queueClient.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving queue properties for %s: %+v", *id, err)
			}

			queueProperties = flattenAccountQueueProperties(queueProps)
		}

		if err := d.Set("queue_properties", queueProperties); err != nil {
			return fmt.Errorf("setting `queue_properties`: %+v", err)
		}

		staticWebsiteProperties := make([]interface{}, 0)
		if supportLevel.supportStaticWebsite {
			accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
			}

			staticWebsiteProperties = flattenAccountStaticWebsiteProperties(staticWebsiteProps)
		}
		if err := d.Set("static_website", staticWebsiteProperties); err != nil {
			return fmt.Errorf("setting `static_website`: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] storage account read for 'blob_properties', 'queue_properties', 'share_properties' and 'static_website' skipped due to DataPlaneAccessOnReadEnabled feature flag being set to 'false'.")
		d.Set("static_website", d.Get("static_website"))
		d.Set("queue_properties", d.Get("queue_properties"))
	}

	return nil
}

func resourceStorageAccountSharePropertiesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// the networking api's only allow a single change to be made to a network layout at once, so let's lock to handle that
	virtualNetworkNames := make([]string, 0)
	if model := existing.Model; model != nil && model.Properties != nil {
		if acls := model.Properties.NetworkAcls; acls != nil {
			if vnr := acls.VirtualNetworkRules; vnr != nil {
				for _, v := range *vnr {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
					if err != nil {
						return err
					}

					networkName := subnetId.VirtualNetworkName
					for _, virtualNetworkName := range virtualNetworkNames {
						if networkName == virtualNetworkName {
							continue
						}
					}
					virtualNetworkNames = append(virtualNetworkNames, networkName)
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// remove this from the cache
	storageClient.RemoveAccountFromCache(*id)

	return nil
}

func expandAccountShareProperties(input []interface{}) fileservice.FileServiceProperties {
	props := fileservice.FileServiceProperties{
		Properties: &fileservice.FileServicePropertiesProperties{
			Cors: &fileservice.CorsRules{
				CorsRules: &[]fileservice.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &fileservice.DeleteRetentionPolicy{
				Enabled: pointer.To(false),
			},
		},
	}

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]interface{})

		props.Properties.ShareDeleteRetentionPolicy = expandAccountShareDeleteRetentionPolicy(v["retention_policy"].([]interface{}))

		props.Properties.Cors = expandAccountSharePropertiesCorsRule(v["cors_rule"].([]interface{}))

		props.Properties.ProtocolSettings = &fileservice.ProtocolSettings{
			Smb: expandAccountSharePropertiesSMB(v["smb"].([]interface{})),
		}
	}

	return props
}

func flattenAccountShareProperties(input *fileservice.FileServiceProperties) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if props := input.Properties; props != nil {
			output = append(output, map[string]interface{}{
				"cors_rule":        flattenAccountSharePropertiesCorsRule(props.Cors),
				"retention_policy": flattenAccountShareDeleteRetentionPolicy(props.ShareDeleteRetentionPolicy),
				"smb":              flattenAccountSharePropertiesSMB(props.ProtocolSettings),
			})
		}
	}

	return output
}

func expandAccountSharePropertiesCorsRule(input []interface{}) *fileservice.CorsRules {
	blobCorsRules := fileservice.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]fileservice.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]interface{})

			allowedMethods := make([]fileservice.AllowedMethods, 0)
			for _, val := range *utils.ExpandStringSlice(item["allowed_methods"].([]interface{})) {
				allowedMethods = append(allowedMethods, fileservice.AllowedMethods(val))
			}
			corsRules = append(corsRules, fileservice.CorsRule{
				AllowedHeaders:  *utils.ExpandStringSlice(item["allowed_headers"].([]interface{})),
				AllowedMethods:  allowedMethods,
				AllowedOrigins:  *utils.ExpandStringSlice(item["allowed_origins"].([]interface{})),
				ExposedHeaders:  *utils.ExpandStringSlice(item["exposed_headers"].([]interface{})),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func flattenAccountSharePropertiesCorsRule(input *fileservice.CorsRules) []interface{} {
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

func expandAccountShareDeleteRetentionPolicy(input []interface{}) *fileservice.DeleteRetentionPolicy {
	result := fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountShareDeleteRetentionPolicy(input *fileservice.DeleteRetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if enabled := input.Enabled; enabled != nil && *enabled {
			days := 0
			if input.Days != nil {
				days = int(*input.Days)
			}

			output = append(output, map[string]interface{}{
				"days": days,
			})
		}
	}

	return output
}

func expandAccountSharePropertiesSMB(input []interface{}) *fileservice.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &fileservice.SmbSetting{
			AuthenticationMethods:    pointer.To(""),
			ChannelEncryption:        pointer.To(""),
			KerberosTicketEncryption: pointer.To(""),
			Versions:                 pointer.To(""),
			Multichannel:             nil,
		}
	}

	v := input[0].(map[string]interface{})

	return &fileservice.SmbSetting{
		AuthenticationMethods:    utils.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        utils.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: utils.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		Versions:                 utils.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		Multichannel: &fileservice.Multichannel{
			Enabled: pointer.To(v["multichannel_enabled"].(bool)),
		},
	}
}

func flattenAccountSharePropertiesSMB(input *fileservice.ProtocolSettings) []interface{} {
	if input == nil || input.Smb == nil {
		return []interface{}{}
	}

	versions := make([]interface{}, 0)
	if input.Smb.Versions != nil {
		versions = utils.FlattenStringSliceWithDelimiter(input.Smb.Versions, ";")
	}

	authenticationMethods := make([]interface{}, 0)
	if input.Smb.AuthenticationMethods != nil {
		authenticationMethods = utils.FlattenStringSliceWithDelimiter(input.Smb.AuthenticationMethods, ";")
	}

	kerberosTicketEncryption := make([]interface{}, 0)
	if input.Smb.KerberosTicketEncryption != nil {
		kerberosTicketEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.KerberosTicketEncryption, ";")
	}

	channelEncryption := make([]interface{}, 0)
	if input.Smb.ChannelEncryption != nil {
		channelEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.ChannelEncryption, ";")
	}

	multichannelEnabled := false
	if input.Smb.Multichannel != nil && input.Smb.Multichannel.Enabled != nil {
		multichannelEnabled = *input.Smb.Multichannel.Enabled
	}

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 && (input.Smb.Multichannel == nil || input.Smb.Multichannel.Enabled == nil) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"authentication_types":            authenticationMethods,
			"channel_encryption_type":         channelEncryption,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"multichannel_enabled":            multichannelEnabled,
			"versions":                        versions,
		},
	}
}
