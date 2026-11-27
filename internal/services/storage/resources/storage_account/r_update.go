// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	tenantId := account.From(meta).TenantId
	storageClient := storageClient2.From(meta).ResourceManager
	client := storageClient.StorageAccounts
	keyVaultClient := keyvaultClient.From(meta)
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	accountTier := storageaccounts.SkuTier(d.Get("account_tier").(string))
	provisionedBillingModelVersion := d.Get("provisioned_billing_model_version").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s%s_%s", accountTier, provisionedBillingModelVersion, replicationType)
	accountKind := storageaccounts.Kind(d.Get("account_kind").(string))

	if accountKind == storageaccounts.KindBlobStorage && storageType == string(storageaccounts.SkuNameStandardZRS) {
		return errors.New("an `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
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
		expandedIdentity, err = identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
	}

	if d.HasChange("access_tier") {
		props.AccessTier = pointer.ToEnum[storageaccounts.AccessTier](d.Get("access_tier").(string))
	}
	if d.HasChange("allowed_copy_scope") {
		props.AllowedCopyScope = pointer.ToEnum[storageaccounts.AllowedCopyScope](d.Get("allowed_copy_scope").(string))
	}
	if d.HasChange("allow_nested_items_to_be_public") {
		props.AllowBlobPublicAccess = pointer.To(d.Get("allow_nested_items_to_be_public").(bool))
	}
	if d.HasChange("cross_tenant_replication_enabled") {
		props.AllowCrossTenantReplication = pointer.To(d.Get("cross_tenant_replication_enabled").(bool))
	}
	if d.HasChange("custom_domain") {
		props.CustomDomain = expandAccountCustomDomain(d.Get("custom_domain").([]any))
	}
	if d.HasChange("customer_managed_key") {
		queueEncryptionKeyType := storageaccounts.KeyType(d.Get("queue_encryption_key_type").(string))
		tableEncryptionKeyType := storageaccounts.KeyType(d.Get("table_encryption_key_type").(string))
		encryptionRaw := d.Get("customer_managed_key").([]any)
		encryption, err := expandAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, encryptionRaw, accountTier, accountKind, *expandedIdentity, queueEncryptionKeyType, tableEncryptionKeyType, provisionedBillingModelVersion)
		if err != nil {
			return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
		}

		// When updating CMK the existing value for `RequireInfrastructureEncryption` gets overwritten which results in
		// an error from the API so we set this back into encryption after it's been overwritten by this update
		existingEnc := existing.Model.Properties.Encryption
		if existingEnc != nil && existingEnc.RequireInfrastructureEncryption != nil {
			encryption.RequireInfrastructureEncryption = existingEnc.RequireInfrastructureEncryption
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

	if d.HasChange("large_file_share_enabled") {
		// largeFileSharesState can only be set to `Enabled` and not `Disabled`, even if it is currently `Disabled`
		if oldValue, newValue := d.GetChange("large_file_share_enabled"); oldValue.(bool) && !newValue.(bool) {
			return errors.New("`large_file_share_enabled` cannot be disabled once it's been enabled")
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
		props.MinimumTlsVersion = pointer.ToEnum[storageaccounts.MinimumTlsVersion](d.Get("min_tls_version").(string))
	}
	if d.HasChange("network_rules") {
		props.NetworkAcls = helpers.ExpandAccountNetworkRules(d.Get("network_rules").([]any), tenantId)
	}

	if d.HasChange("public_network_access") {
		props.PublicNetworkAccess = pointer.ToEnum[storageaccounts.PublicNetworkAccess](d.Get("public_network_access").(string))
	}

	if !features.SixPointOh() {
		if d.HasChange("public_network_access_enabled") {
			publicNetworkAccess := storageaccounts.PublicNetworkAccessDisabled
			if d.Get("public_network_access_enabled").(bool) {
				publicNetworkAccess = storageaccounts.PublicNetworkAccessEnabled
			}
			props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
		}
	}

	if d.HasChange("routing") {
		props.RoutingPreference = expandAccountRoutingPreference(d.Get("routing").([]any))
	}
	if d.HasChange("sas_policy") {
		// TODO: Currently, there is no way to represent a `null` value in the payload - instead it will be omitted, `sas_policy` can not be disabled once enabled.
		props.SasPolicy = expandAccountSASPolicy(d.Get("sas_policy").([]any))
	}
	if d.HasChange("sftp_enabled") {
		props.IsSftpEnabled = pointer.To(d.Get("sftp_enabled").(bool))
	}
	if d.HasChange("immutability_policy") {
		props.ImmutableStorageWithVersioning = expandAccountImmutabilityPolicy(d.Get("immutability_policy").([]any))
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

	// avoid triggering an unnecessary update (`PUT`) API request if there are no other changes besides the below items
	// `account_replication_type` may only require a migration request without the main `PUT` request
	// `azure_files_authentication`, `blob_properties`, and `share_properties` ecah use their own API requests
	updateRequired := d.HasChangesExcept("account_replication_type", "azure_files_authentication", "blob_properties", "share_properties")

	migrationRequired := false
	if d.HasChange("account_replication_type") {
		// certain changes to `account_replication_type` require a storage migration, in this scenario we'll omit the updated SKU and trigger a migration instead
		o, n := d.GetChange("account_replication_type")
		oRT, nRT := o.(string), n.(string)
		if storageReplicationZonalMigrationPairs[oRT] == nRT {
			migrationRequired = true
		}

		if !migrationRequired {
			payload.Sku = storageaccounts.Sku{
				// storageType is derived from "account_replication_type", "account_tier" (force-new) and "provisioned_billing_model_version" (force-new)
				Name: storageaccounts.SkuName(storageType),
			}
			updateRequired = true
		}
	}
	if d.HasChange("identity") {
		payload.Identity = expandedIdentity
	}
	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if updateRequired {
		if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	// azure_files_authentication must be the last to be updated, cause it'll occupy the storage account for several minutes after receiving the response 200 OK. Issue: https://github.com/Azure/azure-rest-api-specs/issues/11272
	if d.HasChange("azure_files_authentication") || d.HasChange("smb_oauth_enabled") {
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

		expandAADFilesAuthentication, err := expandAccountAzureFilesAuthentication(d.Get("azure_files_authentication").([]any), d.Get("smb_oauth_enabled").(bool))
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

	// Following are updates to the sub-services
	supportLevel := helpers.AvailableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if d.HasChange("blob_properties") {
		if !supportLevel.SupportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, d.Get("blob_properties").([]any))
		if err != nil {
			return err
		}

		if blobProperties.Properties.IsVersioningEnabled != nil && *blobProperties.Properties.IsVersioningEnabled && d.Get("is_hns_enabled").(bool) {
			return errors.New("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		// Disable restore_policy first. Disabling restore_policy and while setting delete_retention_policy.allow_permanent_delete to true cause error.
		// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
		if v := d.Get("blob_properties.0.restore_policy"); d.HasChange("blob_properties.0.restore_policy") && len(v.([]any)) == 0 {
			log.Print("[DEBUG] Disabling RestorePolicy prior to changing DeleteRetentionPolicy")
			blobPayload := blobservices.BlobServiceProperties{
				Properties: &blobservices.BlobServicePropertiesProperties{
					RestorePolicy: expandAccountBlobPropertiesRestorePolicy(v.([]any)),
				},
			}
			if _, err := storageClient.BlobServices.SetServiceProperties(ctx, *id, blobPayload); err != nil {
				return fmt.Errorf("updating Azure Storage Account blob restore policy %q: %+v", id.StorageAccountName, err)
			}
		}

		if d.Get("dns_endpoint_type").(string) == string(storageaccounts.DnsEndpointTypeAzureDnsZone) {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		if _, err = storageClient.BlobServices.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties` for %s: %+v", *id, err)
		}
	}

	if d.HasChange("share_properties") {
		if !supportLevel.SupportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		sharePayload := expandAccountShareProperties(d.Get("share_properties").([]any))
		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			if sharePayload.Properties.ProtocolSettings.Smb.Multichannel != nil && sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled != nil {
				if *sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled {
					return errors.New("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}
			}

			sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
		}

		if _, err = storageClient.FileServices.SetServiceProperties(ctx, *id, sharePayload); err != nil {
			return fmt.Errorf("updating File Share Properties for %s: %+v", *id, err)
		}
	}

	if migrationRequired {
		migrationPayload := storageaccounts.StorageAccountMigration{
			Properties: storageaccounts.StorageAccountMigrationProperties{
				TargetSkuName: storageaccounts.SkuName(storageType),
			},
		}

		if _, err := client.CustomerInitiatedMigration(ctx, *id, migrationPayload); err != nil {
			return fmt.Errorf("triggering account migration for %s: %+v", id, err)
		}
		d.Set("account_replication_type_migration_in_progress", true)
		d.Set("account_replication_type_migrating_to", replicationType)
	}

	return resourceStorageAccountRead(d, meta)
}
