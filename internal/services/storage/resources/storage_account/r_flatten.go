// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccountmigrations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceStorageAccountFlatten(ctx context.Context, d *pluginsdk.ResourceData, id commonids.StorageAccountId, account *storageaccounts.StorageAccount, meta any, includeResource bool) error {
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	d.Set("name", id.StorageAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if err := pluginsdk.SetResourceIdentityData(d, pointer.To(id)); err != nil {
		return fmt.Errorf("setting resource identity data: %+v", err)
	}

	if !includeResource {
		return nil
	}

	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	listKeysOpts := storageaccounts.DefaultListKeysOperationOptions()
	listKeysOpts.Expand = pointer.To(storageaccounts.ExpandKerb)

	supportLevel := helpers.StorageAccountServiceSupportLevel{}

	var accountKind storageaccounts.Kind
	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference
	if account.Kind != nil {
		accountKind = *account.Kind
	}
	d.Set("account_kind", string(accountKind))

	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""
	provisionBillingModelVer := ""
	if sku := account.Sku; sku != nil {
		skuNameSegs := strings.Split(string(sku.Name), "_")
		accountTierAndProvisionBillingModelVer := skuNameSegs[0]
		accountReplicationType = skuNameSegs[1]
		if sku.Tier != nil {
			accountTier = *sku.Tier
			provisionBillingModelVer = strings.TrimPrefix(accountTierAndProvisionBillingModelVer, string(accountTier))
		}
	}
	d.Set("account_tier", string(accountTier))
	d.Set("account_replication_type", accountReplicationType)
	d.Set("provisioned_billing_model_version", provisionBillingModelVer)

	d.Set("edge_zone", flattenEdgeZone(account.ExtendedLocation))
	d.Set("location", location.Normalize(account.Location))

	if props := account.Properties; props != nil {
		primaryEndpoints = props.PrimaryEndpoints
		routingPreference = props.RoutingPreference
		secondaryEndpoints = props.SecondaryEndpoints

		d.Set("access_tier", pointer.From(props.AccessTier))
		d.Set("allowed_copy_scope", pointer.From(props.AllowedCopyScope))
		if err := d.Set("azure_files_authentication", flattenAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
			return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
		}
		smbOAuthEnabled := false
		if afa := props.AzureFilesIdentityBasedAuthentication; afa != nil && afa.SmbOAuthSettings != nil {
			smbOAuthEnabled = pointer.From(afa.SmbOAuthSettings.IsSmbOAuthEnabled)
		}
		d.Set("smb_oauth_enabled", smbOAuthEnabled)
		d.Set("cross_tenant_replication_enabled", pointer.From(props.AllowCrossTenantReplication))
		d.Set("https_traffic_only_enabled", pointer.From(props.SupportsHTTPSTrafficOnly))
		d.Set("is_hns_enabled", pointer.From(props.IsHnsEnabled))
		d.Set("nfsv3_enabled", pointer.From(props.IsNfsV3Enabled))
		d.Set("primary_location", pointer.From(props.PrimaryLocation))
		if err := d.Set("routing", flattenAccountRoutingPreference(props.RoutingPreference)); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
		d.Set("secondary_location", pointer.From(props.SecondaryLocation))
		d.Set("sftp_enabled", pointer.From(props.IsSftpEnabled))
		d.Set("allow_nested_items_to_be_public", pointer.From(props.AllowBlobPublicAccess))
		d.Set("default_to_oauth_authentication", pointer.From(props.DefaultToOAuthAuthentication))

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

		if !features.SixPointOh() {
			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != storageaccounts.PublicNetworkAccessEnabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
		}
		d.Set("public_network_access", pointer.FromEnum(props.PublicNetworkAccess))

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
		if err := d.Set("network_rules", helpers.FlattenAccountNetworkRules(props.NetworkAcls)); err != nil {
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

		customerManagedKey, err := flattenAccountCustomerManagedKey(props.Encryption)
		if err != nil {
			return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
		}

		if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		if err := d.Set("sas_policy", flattenAccountSASPolicy(props.SasPolicy)); err != nil {
			return fmt.Errorf("setting `sas_policy`: %+v", err)
		}

		supportLevel = helpers.AvailableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)
	}

	flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(account.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := tags.FlattenAndSet(d, account.Tags); err != nil {
		return err
	}

	endpoints := flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	endpoints.set(d)

	keys, err := client.ListKeys(ctx, id, listKeysOpts)
	if err != nil {
		hasWriteLock := response.WasConflict(keys.HttpResponse)
		doesntHavePermissions := response.WasForbidden(keys.HttpResponse) || response.WasStatusCode(keys.HttpResponse, http.StatusUnauthorized)
		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	storageAccountKeys := make([]storageaccounts.StorageAccountKey, 0)
	if keys.Model != nil && keys.Model.Keys != nil {
		storageAccountKeys = *keys.Model.Keys
	}
	keysAndConnectionStrings := flattenAccountAccessKeysAndConnectionStrings(id.StorageAccountName, *storageDomainSuffix, storageAccountKeys, endpoints)
	keysAndConnectionStrings.set(d)

	blobProperties := make([]any, 0)
	if supportLevel.SupportBlob {
		blobProps, err := storageClient.BlobServices.GetServiceProperties(ctx, id)
		if err != nil {
			return fmt.Errorf("reading blob properties for %s: %+v", id, err)
		}

		blobProperties = flattenAccountBlobServiceProperties(blobProps.Model)
	}
	if err := d.Set("blob_properties", blobProperties); err != nil {
		return fmt.Errorf("setting `blob_properties` for %s: %+v", id, err)
	}

	shareProperties := make([]any, 0)
	if supportLevel.SupportShare {
		shareProps, err := storageClient.FileServices.GetServiceProperties(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving share properties for %s: %+v", id, err)
		}

		shareProperties = flattenAccountShareProperties(shareProps.Model)
	}
	if err := d.Set("share_properties", shareProperties); err != nil {
		return fmt.Errorf("setting `share_properties` for %s: %+v", id, err)
	}

	resp, err := storageClient.StorageAccountMigrations.StorageAccountsGetCustomerInitiatedMigration(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving migration status for %s: %+v", id, err)
	}

	migrationInProgress, migratingTo := false, ""
	if resp.Model != nil {
		props := resp.Model.Properties
		migrationInProgress = pointer.From(props.MigrationStatus) == storageaccountmigrations.MigrationStatusInProgress || pointer.From(props.MigrationStatus) == storageaccountmigrations.MigrationStatusSubmittedForConversion
		if migrationInProgress {
			if splitSKU := strings.Split(string(props.TargetSkuName), "_"); len(splitSKU) == 2 {
				migratingTo = splitSKU[1]
			}
		}
	}
	d.Set("account_replication_type_migration_in_progress", migrationInProgress)
	d.Set("account_replication_type_migrating_to", migratingTo)

	return nil
}

func flattenAccountCustomDomain(input *storageaccounts.CustomDomain) []any {
	output := make([]any, 0)
	if input != nil {
		output = append(output, map[string]any{
			// use_subdomain isn't returned
			"name": input.Name,
		})
	}
	return output
}

func flattenAccountCustomerManagedKey(input *storageaccounts.Encryption) ([]any, error) {
	output := make([]any, 0)

	if input != nil && pointer.From(input.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
		cmk := make(map[string]any)

		if props := input.Identity; props != nil {
			cmk["user_assigned_identity_id"] = pointer.From(props.UserAssignedIdentity)
		}

		if props := input.Keyvaultproperties; props != nil {
			keyID, err := keyvault.NewNestedItemID(pointer.From(props.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(props.Keyname), pointer.From(props.Keyversion))
			if err != nil {
				return output, err
			}
			cmk["key_vault_key_id"] = keyID.ID()
		}

		output = append(output, cmk)
	}

	return output, nil
}

func flattenAccountImmutabilityPolicy(input *storageaccounts.ImmutableStorageAccount) []any {
	if input == nil || input.ImmutabilityPolicy == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"allow_protected_append_writes": input.ImmutabilityPolicy.AllowProtectedAppendWrites,
			"period_since_creation_in_days": input.ImmutabilityPolicy.ImmutabilityPeriodSinceCreationInDays,
			"state":                         input.ImmutabilityPolicy.State,
		},
	}
}

func flattenAccountActiveDirectoryProperties(input *storageaccounts.ActiveDirectoryProperties) []any {
	output := make([]any, 0)
	if input != nil {
		output = append(output, map[string]any{
			"domain_guid":         input.DomainGuid,
			"domain_name":         input.DomainName,
			"domain_sid":          pointer.From(input.DomainSid),
			"forest_name":         pointer.From(input.ForestName),
			"netbios_domain_name": pointer.From(input.NetBiosDomainName),
			"storage_sid":         pointer.From(input.AzureStorageSid),
		})
	}
	return output
}

func flattenAccountAzureFilesAuthentication(input *storageaccounts.AzureFilesIdentityBasedAuthentication) []any {
	if input == nil || input.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsNone {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"active_directory":               flattenAccountActiveDirectoryProperties(input.ActiveDirectoryProperties),
			"directory_type":                 input.DirectoryServiceOptions,
			"default_share_level_permission": input.DefaultSharePermission,
		},
	}
}

func flattenAccountRoutingPreference(input *storageaccounts.RoutingPreference) []any {
	output := make([]any, 0)

	if input != nil {
		routingChoice := ""
		if input.RoutingChoice != nil {
			routingChoice = string(*input.RoutingChoice)
		}

		output = append(output, map[string]any{
			"choice":                      routingChoice,
			"publish_internet_endpoints":  pointer.From(input.PublishInternetEndpoints),
			"publish_microsoft_endpoints": pointer.From(input.PublishMicrosoftEndpoints),
		})
	}

	return output
}

func flattenAccountBlobServiceProperties(input *blobservices.BlobServiceProperties) []any {
	if input == nil || input.Properties == nil {
		return []any{}
	}

	flattenedCorsRules := make([]any, 0)
	if corsRules := input.Properties.Cors; corsRules != nil {
		flattenedCorsRules = flattenAccountBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]any, 0)
	if deletePolicy := input.Properties.DeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenAccountBlobDeleteRetentionPolicy(deletePolicy)
	}

	flattenedRestorePolicy := make([]any, 0)
	if restorePolicy := input.Properties.RestorePolicy; restorePolicy != nil {
		flattenedRestorePolicy = flattenAccountBlobPropertiesRestorePolicy(restorePolicy)
	}

	flattenedContainerDeletePolicy := make([]any, 0)
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

	defaultServiceVersion := pointer.From(input.Properties.DefaultServiceVersion)

	var LastAccessTimeTrackingPolicy bool
	if v := input.Properties.LastAccessTimeTrackingPolicy; v != nil {
		LastAccessTimeTrackingPolicy = v.Enable
	}

	return []any{
		map[string]any{
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

func flattenAccountBlobDeleteRetentionPolicy(input *blobservices.DeleteRetentionPolicy) []any {
	deleteRetentionPolicy := make([]any, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]any{
			"days":                     days,
			"permanent_delete_enabled": pointer.From(input.AllowPermanentDelete),
		})
	}

	return deleteRetentionPolicy
}

func flattenAccountBlobContainerDeleteRetentionPolicy(input *blobservices.DeleteRetentionPolicy) []any {
	deleteRetentionPolicy := make([]any, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]any{
			"days": days,
		})
	}

	return deleteRetentionPolicy
}

func flattenAccountBlobPropertiesRestorePolicy(input *blobservices.RestorePolicyProperties) []any {
	restorePolicy := make([]any, 0)

	if input == nil {
		return restorePolicy
	}

	if enabled := input.Enabled; enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		restorePolicy = append(restorePolicy, map[string]any{
			"days": days,
		})
	}

	return restorePolicy
}

func flattenAccountBlobPropertiesCorsRule(input *blobservices.CorsRules) []any {
	corsRules := make([]any, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]any{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func flattenAccountShareProperties(input *fileservices.FileServiceProperties) []any {
	output := make([]any, 0)

	if input != nil {
		if props := input.Properties; props != nil {
			output = append(output, map[string]any{
				"cors_rule":        flattenAccountSharePropertiesCorsRule(props.Cors),
				"retention_policy": flattenAccountShareDeleteRetentionPolicy(props.ShareDeleteRetentionPolicy),
				"smb":              flattenAccountSharePropertiesSMB(props.ProtocolSettings),
			})
		}
	}

	return output
}

func flattenAccountSharePropertiesCorsRule(input *fileservices.CorsRules) []any {
	corsRules := make([]any, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]any{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func flattenAccountShareDeleteRetentionPolicy(input *fileservices.DeleteRetentionPolicy) []any {
	output := make([]any, 0)

	if input != nil {
		if enabled := input.Enabled; enabled != nil && *enabled {
			days := 0
			if input.Days != nil {
				days = int(*input.Days)
			}

			output = append(output, map[string]any{
				"days": days,
			})
		}
	}

	return output
}

func flattenAccountSharePropertiesSMB(input *fileservices.ProtocolSettings) []any {
	if input == nil || input.Smb == nil {
		return []any{}
	}

	versions := make([]any, 0)
	if input.Smb.Versions != nil {
		versions = pluginsdk.FlattenStringSliceWithDelimiter(input.Smb.Versions, ";")
	}

	authenticationMethods := make([]any, 0)
	if input.Smb.AuthenticationMethods != nil {
		authenticationMethods = pluginsdk.FlattenStringSliceWithDelimiter(input.Smb.AuthenticationMethods, ";")
	}

	kerberosTicketEncryption := make([]any, 0)
	if input.Smb.KerberosTicketEncryption != nil {
		kerberosTicketEncryption = pluginsdk.FlattenStringSliceWithDelimiter(input.Smb.KerberosTicketEncryption, ";")
	}

	channelEncryption := make([]any, 0)
	if input.Smb.ChannelEncryption != nil {
		channelEncryption = pluginsdk.FlattenStringSliceWithDelimiter(input.Smb.ChannelEncryption, ";")
	}

	multichannelEnabled := false
	if input.Smb.Multichannel != nil && input.Smb.Multichannel.Enabled != nil {
		multichannelEnabled = *input.Smb.Multichannel.Enabled
	}

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 && (input.Smb.Multichannel == nil || input.Smb.Multichannel.Enabled == nil) {
		return []any{}
	}

	return []any{
		map[string]any{
			"authentication_types":            authenticationMethods,
			"channel_encryption_type":         channelEncryption,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"multichannel_enabled":            multichannelEnabled,
			"versions":                        versions,
		},
	}
}

func flattenAccountSASPolicy(input *storageaccounts.SasPolicy) []any {
	output := make([]any, 0)

	if input != nil {
		output = append(output, map[string]any{
			"expiration_action": string(input.ExpirationAction),
			"expiration_period": input.SasExpirationPeriod,
		})
	}

	return output
}

func flattenEdgeZone(input *edgezones.Model) string {
	output := ""
	if input != nil {
		output = edgezones.Normalize(input.Name)
	}
	return output
}
