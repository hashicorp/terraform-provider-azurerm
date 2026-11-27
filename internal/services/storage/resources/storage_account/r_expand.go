// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAccountCustomDomain(input []any) *storageaccounts.CustomDomain {
	if len(input) == 0 {
		return &storageaccounts.CustomDomain{}
	}

	domain := input[0].(map[string]any)
	return &storageaccounts.CustomDomain{
		Name:             domain["name"].(string),
		UseSubDomainName: pointer.To(domain["use_subdomain"].(bool)),
	}
}

func expandAccountCustomerManagedKey(ctx context.Context, keyVaultClient *client.Client, subscriptionId string, input []any, accountTier storageaccounts.SkuTier, accountKind storageaccounts.Kind, expandedIdentity identity.LegacySystemAndUserAssignedMap, queueEncryptionKeyType, tableEncryptionKeyType storageaccounts.KeyType, billingModel string) (*storageaccounts.Encryption, error) {
	if !features.SixPointOh() {
		if accountKind == storageaccounts.KindStorage {
			if queueEncryptionKeyType == storageaccounts.KeyTypeAccount {
				return nil, fmt.Errorf("`queue_encryption_key_type = %q` cannot be used with account kind `%q`", string(storageaccounts.KeyTypeAccount), string(storageaccounts.KindStorage))
			}
			if tableEncryptionKeyType == storageaccounts.KeyTypeAccount {
				return nil, fmt.Errorf("`table_encryption_key_type = %q` cannot be used with account kind `%q`", string(storageaccounts.KeyTypeAccount), string(storageaccounts.KindStorage))
			}
		}
	}

	if len(input) == 0 {
		return &storageaccounts.Encryption{
			KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointStorage),
			Services: &storageaccounts.EncryptionServices{
				Queue: &storageaccounts.EncryptionService{
					KeyType: pointer.To(queueEncryptionKeyType),
				},
				Table: &storageaccounts.EncryptionService{
					KeyType: pointer.To(tableEncryptionKeyType),
				},
			},
		}, nil
	}

	if accountTier != storageaccounts.SkuTierPremium && accountKind != storageaccounts.KindStorageVTwo && (accountKind != storageaccounts.KindFileStorage || billingModel != "V2") {
		return nil, errors.New("customer managed key can only be used with account kind `StorageV2`, account tier `Premium`, or account kind `FileStorage` with billing model `V2`")
	}

	if expandedIdentity.Type != identity.TypeUserAssigned && expandedIdentity.Type != identity.TypeSystemAssignedUserAssigned {
		return nil, fmt.Errorf("customer managed key can only be configured when the storage account uses a `UserAssigned` or `SystemAssigned, UserAssigned` managed identity but got %q", string(expandedIdentity.Type))
	}

	v := input[0].(map[string]any)

	keyID, err := keyvault.ParseNestedItemID(v["key_vault_key_id"].(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	if !keyID.IsManagedHSM() {
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
		keyVaultIdRaw, err := keyVaultClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyID.KeyVaultBaseURL)
		if err != nil {
			return nil, err
		}
		if keyVaultIdRaw == nil {
			return nil, fmt.Errorf("unable to find the Resource Manager ID for the Key Vault URI %q in %s", keyID.KeyVaultBaseURL, subscriptionResourceId)
		}
		keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
		if err != nil {
			return nil, err
		}

		vaultsClient := keyVaultClient.VaultsClient
		keyVault, err := vaultsClient.Get(ctx, *keyVaultId)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *keyVaultId, err)
		}

		softDeleteEnabled := false
		purgeProtectionEnabled := false
		if model := keyVault.Model; model != nil {
			if esd := model.Properties.EnableSoftDelete; esd != nil {
				softDeleteEnabled = *esd
			}
			if epp := model.Properties.EnablePurgeProtection; epp != nil {
				purgeProtectionEnabled = *epp
			}
		}
		if !softDeleteEnabled || !purgeProtectionEnabled {
			return nil, fmt.Errorf("%s must be configured for both Purge Protection and Soft Delete", *keyVaultId)
		}
	}

	return &storageaccounts.Encryption{
		Services: &storageaccounts.EncryptionServices{
			Blob: &storageaccounts.EncryptionService{
				Enabled: pointer.To(true),
				KeyType: pointer.To(storageaccounts.KeyTypeAccount),
			},
			File: &storageaccounts.EncryptionService{
				Enabled: pointer.To(true),
				KeyType: pointer.To(storageaccounts.KeyTypeAccount),
			},
			Queue: &storageaccounts.EncryptionService{
				KeyType: pointer.To(queueEncryptionKeyType),
			},
			Table: &storageaccounts.EncryptionService{
				KeyType: pointer.To(tableEncryptionKeyType),
			},
		},
		Identity: &storageaccounts.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(v["user_assigned_identity_id"].(string)),
		},
		KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointKeyvault),
		Keyvaultproperties: &storageaccounts.KeyVaultProperties{
			Keyname:     pointer.To(keyID.Name),
			Keyversion:  pointer.To(keyID.Version),
			Keyvaulturi: pointer.To(keyID.KeyVaultBaseURL),
		},
	}, nil
}

func expandAccountImmutabilityPolicy(input []any) *storageaccounts.ImmutableStorageAccount {
	if len(input) == 0 {
		return &storageaccounts.ImmutableStorageAccount{}
	}

	v := input[0].(map[string]any)
	return &storageaccounts.ImmutableStorageAccount{
		Enabled: pointer.To(true),
		ImmutabilityPolicy: &storageaccounts.AccountImmutabilityPolicyProperties{
			AllowProtectedAppendWrites:            pointer.To(v["allow_protected_append_writes"].(bool)),
			ImmutabilityPeriodSinceCreationInDays: pointer.To(int64(v["period_since_creation_in_days"].(int))),
			State:                                 pointer.ToEnum[storageaccounts.AccountImmutabilityPolicyState](v["state"].(string)),
		},
	}
}

func expandAccountActiveDirectoryProperties(input []any) *storageaccounts.ActiveDirectoryProperties {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]any)

	output := &storageaccounts.ActiveDirectoryProperties{
		DomainGuid: pointer.To(m["domain_guid"].(string)),
		DomainName: pointer.To(m["domain_name"].(string)),
	}
	if v := m["storage_sid"]; v != "" {
		output.AzureStorageSid = pointer.To(v.(string))
	}
	if v := m["domain_sid"]; v != "" {
		output.DomainSid = pointer.To(v.(string))
	}
	if v := m["forest_name"]; v != "" {
		output.ForestName = pointer.To(v.(string))
	}
	if v := m["netbios_domain_name"]; v != "" {
		output.NetBiosDomainName = pointer.To(v.(string))
	}
	return output
}

func expandAccountAzureFilesAuthentication(input []any, smbOAuthEnabled bool) (*storageaccounts.AzureFilesIdentityBasedAuthentication, error) {
	smbOAuthSettings := &storageaccounts.SmbOAuthSettings{
		IsSmbOAuthEnabled: pointer.To(smbOAuthEnabled),
	}

	if len(input) == 0 {
		return &storageaccounts.AzureFilesIdentityBasedAuthentication{
			DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsNone,
			SmbOAuthSettings:        smbOAuthSettings,
		}, nil
	}

	v := input[0].(map[string]any)
	output := storageaccounts.AzureFilesIdentityBasedAuthentication{
		DirectoryServiceOptions: storageaccounts.DirectoryServiceOptions(v["directory_type"].(string)),
		SmbOAuthSettings:        smbOAuthSettings,
	}
	if output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAD ||
		output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAADDS ||
		output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAADKERB {
		ad := expandAccountActiveDirectoryProperties(v["active_directory"].([]any))

		if output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAD {
			if ad == nil {
				return nil, errors.New("`active_directory` is required when `directory_type` is `AD`")
			}
			if ad.AzureStorageSid == nil {
				return nil, errors.New("`active_directory.0.storage_sid` is required when `directory_type` is `AD`")
			}
			if ad.DomainSid == nil {
				return nil, errors.New("`active_directory.0.domain_sid` is required when `directory_type` is `AD`")
			}
			if ad.ForestName == nil {
				return nil, errors.New("`active_directory.0.forest_name` is required when `directory_type` is `AD`")
			}
			if ad.NetBiosDomainName == nil {
				return nil, errors.New("`active_directory.0.netbios_domain_name` is required when `directory_type` is `AD`")
			}
		}

		output.ActiveDirectoryProperties = ad
		output.DefaultSharePermission = pointer.ToEnum[storageaccounts.DefaultSharePermission](v["default_share_level_permission"].(string))
	}

	return &output, nil
}

func expandAccountRoutingPreference(input []any) *storageaccounts.RoutingPreference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &storageaccounts.RoutingPreference{
		PublishMicrosoftEndpoints: pointer.To(v["publish_microsoft_endpoints"].(bool)),
		PublishInternetEndpoints:  pointer.To(v["publish_internet_endpoints"].(bool)),
		RoutingChoice:             pointer.ToEnum[storageaccounts.RoutingChoice](v["choice"].(string)),
	}
}

func expandAccountBlobServiceProperties(kind storageaccounts.Kind, input []any) (*blobservices.BlobServiceProperties, error) {
	props := blobservices.BlobServicePropertiesProperties{
		ChangeFeed: &blobservices.ChangeFeed{
			Enabled: pointer.To(false),
		},
		Cors: &blobservices.CorsRules{
			CorsRules: &[]blobservices.CorsRule{},
		},
		DeleteRetentionPolicy: &blobservices.DeleteRetentionPolicy{
			Enabled: pointer.To(false),
		},
		IsVersioningEnabled:          pointer.To(false),
		LastAccessTimeTrackingPolicy: &blobservices.LastAccessTimeTrackingPolicy{},
	}

	if !features.SixPointOh() {
		// `Storage` (v1) kind doesn't support:
		// - LastAccessTimeTrackingPolicy: Confirmed by SRP.
		// - ChangeFeed: See https://learn.microsoft.com/azure/storage/blobs/storage-blob-change-feed?tabs=azure-portal#enable-and-disable-the-change-feed.
		// - Versioning: See https://learn.microsoft.com/azure/storage/blobs/versioning-overview#how-blob-versioning-works
		// - Restore Policy: See https://learn.microsoft.com/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
		if kind == storageaccounts.KindStorage {
			props.LastAccessTimeTrackingPolicy = nil
			props.ChangeFeed = nil
			props.IsVersioningEnabled = nil
		}
	}

	if len(input) > 0 {
		v := input[0].(map[string]any)

		deletePolicyRaw := v["delete_retention_policy"].([]any)
		props.DeleteRetentionPolicy = expandAccountBlobDeleteRetentionPolicy(deletePolicyRaw)

		containerDeletePolicyRaw := v["container_delete_retention_policy"].([]any)
		props.ContainerDeleteRetentionPolicy = expandAccountBlobContainerDeleteRetentionPolicy(containerDeletePolicyRaw)

		corsRaw := v["cors_rule"].([]any)
		props.Cors = expandAccountBlobPropertiesCors(corsRaw)

		props.IsVersioningEnabled = pointer.To(v["versioning_enabled"].(bool))

		if version, ok := v["default_service_version"].(string); ok && version != "" {
			props.DefaultServiceVersion = pointer.To(version)
		}

		lastAccessTimeEnabled := v["last_access_time_enabled"].(bool)
		changeFeedEnabled := v["change_feed_enabled"].(bool)
		changeFeedRetentionInDays := v["change_feed_retention_in_days"].(int)
		restorePolicyRaw := v["restore_policy"].([]any)
		versioningEnabled := v["versioning_enabled"].(bool)
		if !features.SixPointOh() && kind == storageaccounts.KindStorage {
			// `Storage` (v1) kind doesn't support:
			// - LastAccessTimeTrackingPolicy
			// - ChangeFeed
			// - Versioning
			// - RestorePolicy
			if lastAccessTimeEnabled {
				return nil, errors.New("`last_access_time_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedEnabled {
				return nil, errors.New("`change_feed_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedRetentionInDays != 0 {
				return nil, errors.New("`change_feed_retention_in_days` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if len(restorePolicyRaw) != 0 {
				return nil, errors.New("`restore_policy` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if versioningEnabled {
				return nil, errors.New("`versioning_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
		} else {
			props.LastAccessTimeTrackingPolicy = &blobservices.LastAccessTimeTrackingPolicy{
				Enable: lastAccessTimeEnabled,
			}
			props.ChangeFeed = &blobservices.ChangeFeed{
				Enabled: pointer.To(changeFeedEnabled),
			}
			if changeFeedRetentionInDays != 0 {
				props.ChangeFeed.RetentionInDays = pointer.To(int64(changeFeedRetentionInDays))
			}
			props.RestorePolicy = expandAccountBlobPropertiesRestorePolicy(restorePolicyRaw)
			props.IsVersioningEnabled = &versioningEnabled
		}

		// Sanity check for the prerequisites of restore_policy
		// Ref: https://learn.microsoft.com/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
		if p := props.RestorePolicy; p != nil && p.Enabled {
			if props.ChangeFeed == nil || props.ChangeFeed.Enabled == nil || !*props.ChangeFeed.Enabled {
				return nil, errors.New("`change_feed_enabled` must be `true` when `restore_policy` is set")
			}
			if props.IsVersioningEnabled == nil || !*props.IsVersioningEnabled {
				return nil, errors.New("`versioning_enabled` must be `true` when `restore_policy` is set")
			}
		}
	}

	return &blobservices.BlobServiceProperties{
		Properties: &props,
	}, nil
}

func expandAccountBlobDeleteRetentionPolicy(input []any) *blobservices.DeleteRetentionPolicy {
	result := blobservices.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]any)

	return &blobservices.DeleteRetentionPolicy{
		Enabled:              pointer.To(true),
		AllowPermanentDelete: pointer.To(policy["permanent_delete_enabled"].(bool)),
		Days:                 pointer.To(int64(policy["days"].(int))),
	}
}

func expandAccountBlobContainerDeleteRetentionPolicy(input []any) *blobservices.DeleteRetentionPolicy {
	result := blobservices.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]any)

	return &blobservices.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func expandAccountBlobPropertiesRestorePolicy(input []any) *blobservices.RestorePolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return pointer.To(blobservices.RestorePolicyProperties{})
	}

	policy := input[0].(map[string]any)

	return &blobservices.RestorePolicyProperties{
		Enabled: true,
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func expandAccountBlobPropertiesCors(input []any) *blobservices.CorsRules {
	blobCorsRules := blobservices.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]blobservices.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]any)

			allowedMethods := make([]blobservices.AllowedMethods, 0)
			for _, val := range *pluginsdk.ExpandStringSlice(item["allowed_methods"].([]any)) {
				allowedMethods = append(allowedMethods, blobservices.AllowedMethods(val))
			}
			corsRules = append(corsRules, blobservices.CorsRule{
				AllowedHeaders:  *pluginsdk.ExpandStringSlice(item["allowed_headers"].([]any)),
				AllowedOrigins:  *pluginsdk.ExpandStringSlice(item["allowed_origins"].([]any)),
				AllowedMethods:  allowedMethods,
				ExposedHeaders:  *pluginsdk.ExpandStringSlice(item["exposed_headers"].([]any)),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func expandAccountShareProperties(input []any) fileservices.FileServiceProperties {
	props := fileservices.FileServiceProperties{
		Properties: &fileservices.FileServicePropertiesProperties{
			Cors: &fileservices.CorsRules{
				CorsRules: &[]fileservices.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &fileservices.DeleteRetentionPolicy{
				Enabled: pointer.To(false),
			},
		},
	}

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]any)

		props.Properties.ShareDeleteRetentionPolicy = expandAccountShareDeleteRetentionPolicy(v["retention_policy"].([]any))

		props.Properties.Cors = expandAccountSharePropertiesCorsRule(v["cors_rule"].([]any))

		props.Properties.ProtocolSettings = &fileservices.ProtocolSettings{
			Smb: expandAccountSharePropertiesSMB(v["smb"].([]any)),
		}
	}

	return props
}

func expandAccountSharePropertiesCorsRule(input []any) *fileservices.CorsRules {
	blobCorsRules := fileservices.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]fileservices.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]any)

			allowedMethods := make([]fileservices.AllowedMethods, 0)
			for _, val := range *pluginsdk.ExpandStringSlice(item["allowed_methods"].([]any)) {
				allowedMethods = append(allowedMethods, fileservices.AllowedMethods(val))
			}
			corsRules = append(corsRules, fileservices.CorsRule{
				AllowedHeaders:  *pluginsdk.ExpandStringSlice(item["allowed_headers"].([]any)),
				AllowedMethods:  allowedMethods,
				AllowedOrigins:  *pluginsdk.ExpandStringSlice(item["allowed_origins"].([]any)),
				ExposedHeaders:  *pluginsdk.ExpandStringSlice(item["exposed_headers"].([]any)),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func expandAccountShareDeleteRetentionPolicy(input []any) *fileservices.DeleteRetentionPolicy {
	result := fileservices.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]any)

	return &fileservices.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func expandAccountSharePropertiesSMB(input []any) *fileservices.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &fileservices.SmbSetting{
			AuthenticationMethods:    pointer.To(""),
			ChannelEncryption:        pointer.To(""),
			KerberosTicketEncryption: pointer.To(""),
			Versions:                 pointer.To(""),
		}
	}

	v := input[0].(map[string]any)

	return &fileservices.SmbSetting{
		AuthenticationMethods:    pluginsdk.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        pluginsdk.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: pluginsdk.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		Versions:                 pluginsdk.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		Multichannel: &fileservices.Multichannel{
			Enabled: pointer.To(v["multichannel_enabled"].(bool)),
		},
	}
}

func expandAccountSASPolicy(input []any) *storageaccounts.SasPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]any)
	return &storageaccounts.SasPolicy{
		ExpirationAction:    storageaccounts.ExpirationAction(raw["expiration_action"].(string)),
		SasExpirationPeriod: raw["expiration_period"].(string),
	}
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}
