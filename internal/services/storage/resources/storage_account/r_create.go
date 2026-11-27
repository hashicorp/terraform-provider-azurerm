// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	tenantId := clientsAccount.From(meta).TenantId
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	storageUtils := storageClient2.From(meta)
	storageClient := storageClient2.From(meta).ResourceManager
	client := storageClient.StorageAccounts
	keyVaultClient := keyvaultClient.From(meta)
	ctx, cancel := timeouts.ForCreate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_account", id.ID())
		}
	}

	accountKind := storageaccounts.Kind(d.Get("account_kind").(string))
	accountTier := storageaccounts.SkuTier(d.Get("account_tier").(string))
	provisionedBillingModelVersion := d.Get("provisioned_billing_model_version").(string)
	replicationType := d.Get("account_replication_type").(string)

	publicNetworkAccess := storageaccounts.PublicNetworkAccessEnabled
	if !pluginsdk.IsExplicitlyNullInConfig(d, "public_network_access") {
		// (features.SixPointOh) In 6.0, the only line needed is the below, it can be inlined in the payload instantiation
		publicNetworkAccess = storageaccounts.PublicNetworkAccess(d.Get("public_network_access").(string))
	}

	if !features.SixPointOh() {
		if !pluginsdk.IsExplicitlyNullInConfig(d, "public_network_access_enabled") && !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = storageaccounts.PublicNetworkAccessDisabled
		}
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	httpsTrafficOnlyEnabled := true
	//nolint:staticcheck
	if v, ok := d.GetOkExists("https_traffic_only_enabled"); ok {
		httpsTrafficOnlyEnabled = v.(bool)
	}

	dnsEndpointType := d.Get("dns_endpoint_type").(string)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)
	nfsV3Enabled := d.Get("nfsv3_enabled").(bool)
	payload := storageaccounts.StorageAccountCreateParameters{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Kind:             accountKind,
		Identity:         expandedIdentity,
		Location:         location.Normalize(d.Get("location").(string)),
		Properties: &storageaccounts.StorageAccountPropertiesCreateParameters{
			AllowBlobPublicAccess:        pointer.To(d.Get("allow_nested_items_to_be_public").(bool)),
			AllowCrossTenantReplication:  pointer.To(d.Get("cross_tenant_replication_enabled").(bool)),
			AllowSharedKeyAccess:         pointer.To(d.Get("shared_access_key_enabled").(bool)),
			DnsEndpointType:              pointer.ToEnum[storageaccounts.DnsEndpointType](dnsEndpointType),
			DefaultToOAuthAuthentication: pointer.To(d.Get("default_to_oauth_authentication").(bool)),
			SupportsHTTPSTrafficOnly:     pointer.To(httpsTrafficOnlyEnabled),
			IsNfsV3Enabled:               pointer.To(nfsV3Enabled),
			IsHnsEnabled:                 pointer.To(isHnsEnabled),
			IsLocalUserEnabled:           pointer.To(d.Get("local_user_enabled").(bool)),
			IsSftpEnabled:                pointer.To(d.Get("sftp_enabled").(bool)),
			MinimumTlsVersion:            pointer.ToEnum[storageaccounts.MinimumTlsVersion](d.Get("min_tls_version").(string)),
			NetworkAcls:                  helpers.ExpandAccountNetworkRules(d.Get("network_rules").([]any), tenantId),
			PublicNetworkAccess:          pointer.To(publicNetworkAccess),
			SasPolicy:                    expandAccountSASPolicy(d.Get("sas_policy").([]any)),
		},
		Sku: storageaccounts.Sku{
			Name: storageaccounts.SkuName(fmt.Sprintf("%s%s_%s", string(accountTier), provisionedBillingModelVersion, replicationType)),
			Tier: pointer.To(accountTier),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v := d.Get("allowed_copy_scope").(string); v != "" {
		payload.Properties.AllowedCopyScope = pointer.ToEnum[storageaccounts.AllowedCopyScope](v)
	}
	azureFilesAuthenticationRaw := d.Get("azure_files_authentication").([]any)
	smbOAuthEnabled := d.Get("smb_oauth_enabled").(bool)
	if len(azureFilesAuthenticationRaw) > 0 || smbOAuthEnabled {
		expandAADFilesAuthentication, err := expandAccountAzureFilesAuthentication(azureFilesAuthenticationRaw, smbOAuthEnabled)
		if err != nil {
			return fmt.Errorf("parsing `azure_files_authentication`: %v", err)
		}
		payload.Properties.AzureFilesIdentityBasedAuthentication = expandAADFilesAuthentication
	}
	if _, ok := d.GetOk("custom_domain"); ok {
		payload.Properties.CustomDomain = expandAccountCustomDomain(d.Get("custom_domain").([]any))
	}
	if v, ok := d.GetOk("immutability_policy"); ok {
		payload.Properties.ImmutableStorageWithVersioning = expandAccountImmutabilityPolicy(v.([]any))
	}

	// BlobStorage does not support ZRS
	if accountKind == storageaccounts.KindBlobStorage && string(payload.Sku.Name) == string(storageaccounts.SkuNameStandardZRS) {
		return errors.New("`account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
	}

	accessTier, accessTierSetInConfig := d.GetOk("access_tier")
	_, skuTierSupported := storageKindsSupportsSkuTier[accountKind]
	if !skuTierSupported && accessTierSetInConfig {
		keys := sortedKeysFromSlice(storageKindsSupportsSkuTier)
		return fmt.Errorf("`access_tier` is only available for accounts of kind set to one of: %+v", strings.Join(keys, " / "))
	}
	if skuTierSupported {
		if !accessTierSetInConfig {
			// default to "Hot"
			accessTier = string(storageaccounts.AccessTierHot)
		}
		payload.Properties.AccessTier = pointer.ToEnum[storageaccounts.AccessTier](accessTier.(string))
	}

	// NFSv3 is supported for standard general-purpose v2 storage accounts and for premium block blob storage accounts.
	// (https://learn.microsoft.com/azure/storage/blobs/network-file-system-protocol-support-how-to#step-5-create-and-configure-a-storage-account)
	if nfsV3Enabled {
		if !isHnsEnabled {
			return errors.New("`nfsv3_enabled` can only be used when `is_hns_enabled` is `true`")
		}

		isPremiumTierAndBlockBlobStorageKind := accountTier == storageaccounts.SkuTierPremium && accountKind == storageaccounts.KindBlockBlobStorage
		isStandardTierAndStorageV2Kind := accountTier == storageaccounts.SkuTierStandard && accountKind == storageaccounts.KindStorageVTwo
		if !isPremiumTierAndBlockBlobStorageKind && !isStandardTierAndStorageV2Kind {
			return errors.New("`nfsv3_enabled` can only be used with account tier `Standard` and account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
		}
	}

	//nolint:staticcheck
	if v, ok := d.GetOkExists("large_file_share_enabled"); ok {
		// @tombuildsstuff: we can't set this to `false` because the API returns:
		//
		// performing Create: unexpected status 400 (400 Bad Request) with error: InvalidRequestPropertyValue: The
		// value 'Disabled' is not allowed for property largeFileSharesState. For more information, see -
		// https://aka.ms/storageaccountlargefilesharestate
		if v.(bool) {
			if _, ok := storageKindsSupportLargeFileShares[accountKind]; !ok {
				keys := sortedKeysFromSlice(storageKindsSupportLargeFileShares)
				return fmt.Errorf("`large_file_shares_enabled` can only be set to `true` with `account_kind` set to one of: %+v", strings.Join(keys, " / "))
			}
			payload.Properties.LargeFileSharesState = pointer.To(storageaccounts.LargeFileSharesStateEnabled)
		}
	}

	if v, ok := d.GetOk("routing"); ok {
		payload.Properties.RoutingPreference = expandAccountRoutingPreference(v.([]any))
	}

	// TODO look into standardizing this across resources that support CMK and at the very least look at improving the UX
	// for encryption of blob, file, table and queue
	//
	// By default (by leaving empty), the table and queue encryption key type is set to "Service". While users can change it to "Account" so that
	// they can further use CMK to encrypt table/queue data. Only the StorageV2 account kind supports the Account key type.
	// Also noted that the blob and file are always using the "Account" key type.
	// See: https://learn.microsoft.com/azure/storage/common/account-encryption-key-create?tabs=portal
	queueEncryptionKeyType := storageaccounts.KeyType(d.Get("queue_encryption_key_type").(string))
	tableEncryptionKeyType := storageaccounts.KeyType(d.Get("table_encryption_key_type").(string))
	encryptionRaw := d.Get("customer_managed_key").([]any)
	encryption, err := expandAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, encryptionRaw, accountTier, accountKind, *expandedIdentity, queueEncryptionKeyType, tableEncryptionKeyType, provisionedBillingModelVersion)
	if err != nil {
		return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
	}

	infrastructureEncryption := d.Get("infrastructure_encryption_enabled").(bool)

	if infrastructureEncryption {
		validPremiumConfiguration := accountTier == storageaccounts.SkuTierPremium && (accountKind == storageaccounts.KindBlockBlobStorage) || accountKind == storageaccounts.KindFileStorage
		validV2Configuration := accountKind == storageaccounts.KindStorageVTwo
		if !validPremiumConfiguration && !validV2Configuration {
			return errors.New("`infrastructure_encryption_enabled` can only be used with account kind `StorageV2`, or account tier `Premium` and account kind is one of `BlockBlobStorage` or `FileStorage`")
		}
		encryption.RequireInfrastructureEncryption = &infrastructureEncryption
	}

	payload.Properties.Encryption = encryption

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// populate the cache
	account, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if account.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if err := storageUtils.AddToCache(id, *account.Model); err != nil {
		return fmt.Errorf("populating cache for %s: %+v", id, err)
	}

	supportLevel := helpers.AvailableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if val, ok := d.GetOk("blob_properties"); ok {
		if !supportLevel.SupportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, val.([]any))
		if err != nil {
			return err
		}

		// See: https://learn.microsoft.com/azure/storage/blobs/versioning-overview#:~:text=Storage%20accounts%20with%20a%20hierarchical%20namespace%20enabled%20for%20use%20with%20Azure%20Data%20Lake%20Storage%20Gen2%20are%20not%20currently%20supported.
		isVersioningEnabled := pointer.From(blobProperties.Properties.IsVersioningEnabled)
		if isVersioningEnabled && isHnsEnabled {
			return errors.New("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		if !isVersioningEnabled {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Conflicting feature 'restorePolicy' is enabled. Please disable it and retry."
				return errors.New("`blob_properties.restore_policy` can't be set when `versioning_enabled` is false")
			}

			immutableStorageWithVersioningEnabled := false
			if props := account.Model.Properties; props != nil {
				if versioning := props.ImmutableStorageWithVersioning; versioning != nil {
					if versioning.ImmutabilityPolicy != nil && versioning.Enabled != nil {
						immutableStorageWithVersioningEnabled = *versioning.Enabled
					}
				}
			}
			if immutableStorageWithVersioningEnabled {
				// Otherwise, API returns: "Conflicting feature 'Account level WORM' is enabled. Please disable it and retry."
				// See: https://learn.microsoft.com/azure/storage/blobs/immutable-policy-configure-version-scope?tabs=azure-portal#prerequisites
				return errors.New("`immutability_policy` can't be set when `versioning_enabled` is false")
			}
		}

		// TODO: This is a temporary limitation on Storage service. Remove this check once the API supports this scenario.
		// See https://github.com/hashicorp/terraform-provider-azurerm/pull/25450#discussion_r1542471667 for the context.
		if dnsEndpointType == string(storageaccounts.DnsEndpointTypeAzureDnsZone) {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		if _, err = storageClient.BlobServices.SetServiceProperties(ctx, id, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties`: %+v", err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		if !supportLevel.SupportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		sharePayload := expandAccountShareProperties(val.([]any))

		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium && sharePayload.Properties != nil && sharePayload.Properties.ProtocolSettings != nil {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			smb := sharePayload.Properties.ProtocolSettings.Smb
			if smb != nil && smb.Multichannel != nil {
				if smb.Multichannel.Enabled != nil && *smb.Multichannel.Enabled {
					return errors.New("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}

				sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
			}
		}

		if _, err = storageClient.FileServices.SetServiceProperties(ctx, id, sharePayload); err != nil {
			return fmt.Errorf("updating `share_properties`: %+v", err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}
