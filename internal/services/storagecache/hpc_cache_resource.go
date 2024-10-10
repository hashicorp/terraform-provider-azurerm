// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/caches"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHPCCache() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHPCCacheCreateOrUpdate,
		Update: resourceHPCCacheCreateOrUpdate,
		Read:   resourceHPCCacheRead,
		Delete: resourceHPCCacheDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := caches.ParseCacheID(id)
			return err
		}),

		Schema: resourceHPCCacheSchema(),
	}
}

func resourceHPCCacheCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure HPC Cache creation.")
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := caches.NewCacheID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing HPC Cache %q: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hpc_cache", id.ID())
		}
	}

	location := d.Get("location").(string)
	cacheSize := d.Get("cache_size_in_gb").(int)
	subnet := d.Get("subnet_id").(string)
	skuName := d.Get("sku_name").(string)

	// SKU Cache Combo Validation
	switch {
	case skuName == "Standard_L4_5G" && cacheSize != 21623:
		return fmt.Errorf("The Standard_L4_5G SKU only supports a cache size of 21623")
	case skuName == "Standard_L9G" && cacheSize != 43246:
		return fmt.Errorf("The Standard_L9G SKU only supports a cache size of 43246")
	case skuName == "Standard_L16G" && cacheSize != 86491:
		return fmt.Errorf("The Standard_L16G SKU only supports a cache size of 86491")
	case (cacheSize == 21623 || cacheSize == 43246 || cacheSize == 86491) && (skuName == "Standard_2G" || skuName == "Standard_4G" || skuName == "Standard_8G"):
		return fmt.Errorf("Incompatible cache size chosen. 21623, 43246 and 86491 are reserved for Read Only resources.")
	}

	var accessPolicies []caches.NfsAccessPolicy
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving existing HPC Cache %q: %v", id, err)
		}
		if model := existing.Model; model != nil {
			if prop := model.Properties; prop != nil {
				if settings := prop.SecuritySettings; settings != nil {
					if policies := settings.AccessPolicies; policies != nil {
						accessPolicies = *policies
					}
				}
			}
		}
	}
	defaultAccessPolicy := expandStorageCacheDefaultAccessPolicy(d.Get("default_access_policy").([]interface{}))
	if defaultAccessPolicy != nil {
		var err error
		accessPolicies, err = CacheInsertOrUpdateAccessPolicy(accessPolicies, *defaultAccessPolicy)
		if err != nil {
			return err
		}
	}

	directorySetting := expandStorageCacheDirectorySettings(d)

	i, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	cache := caches.Cache{
		Name:     pointer.To(name),
		Location: pointer.To(location),
		Properties: &caches.CacheProperties{
			CacheSizeGB:     utils.Int64(int64(cacheSize)),
			Subnet:          pointer.To(subnet),
			NetworkSettings: expandStorageCacheNetworkSettings(d),
			SecuritySettings: &caches.CacheSecuritySettings{
				AccessPolicies: &accessPolicies,
			},
			DirectoryServicesSettings: directorySetting,
		},
		Sku: &caches.CacheSku{
			Name: pointer.To(skuName),
		},
		Identity: i,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if !d.IsNewResource() {
		oldKeyVaultKeyId, newKeyVaultKeyId := d.GetChange("key_vault_key_id")
		if (oldKeyVaultKeyId.(string) != "" && newKeyVaultKeyId.(string) == "") || (oldKeyVaultKeyId.(string) == "" && newKeyVaultKeyId.(string) != "") {
			return fmt.Errorf("`key_vault_key_id` can not be added or removed after HPC Cache is created")
		}
	}

	requireAdditionalUpdate := false
	if v, ok := d.GetOk("key_vault_key_id"); ok {
		autoKeyRotationEnabled := d.Get("automatically_rotate_key_to_latest_enabled").(bool)
		if !d.IsNewResource() && d.HasChange("key_vault_key_id") && autoKeyRotationEnabled {
			// It is by design that `automatically_rotate_key_to_latest_enabled` changes to `false` when `key_vault_key_id` is changed, needs to do an additional update to set it back
			requireAdditionalUpdate = true
		}
		// For new created resource `automatically_rotate_key_to_latest_enabled` needs an additional update to set it to true to.
		if d.IsNewResource() && autoKeyRotationEnabled {
			requireAdditionalUpdate = true
		}

		keyVaultKeyId := v.(string)
		keyVaultDetails, err := storageCacheRetrieveKeyVault(ctx, keyVaultsClient, subscriptionId, keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("validating Key Vault Key %q for HPC Cache: %+v", keyVaultKeyId, err)
		}
		if azure.NormalizeLocation(keyVaultDetails.location) != azure.NormalizeLocation(location) {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for HPC Cache: Key Vault must be in the same region as HPC Cache!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}
		if !keyVaultDetails.softDeleteEnabled {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for HPC Cache: Soft Delete must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}
		if !keyVaultDetails.purgeProtectionEnabled {
			return fmt.Errorf("validating Key Vault %q (Resource Group %q) for HPC Cache: Purge Protection must be enabled but it isn't!", keyVaultDetails.keyVaultName, keyVaultDetails.resourceGroupName)
		}

		cache.Properties.EncryptionSettings = &caches.CacheEncryptionSettings{
			KeyEncryptionKey: &caches.KeyVaultKeyReference{
				KeyURL: keyVaultKeyId,
				SourceVault: caches.KeyVaultKeyReferenceSourceVault{
					Id: pointer.To(keyVaultDetails.keyVaultId),
				},
			},
			RotationToLatestKeyVersionEnabled: pointer.To(autoKeyRotationEnabled),
		}
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, cache); err != nil {
		return fmt.Errorf("creating/updating HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if requireAdditionalUpdate {
		if err := client.CreateOrUpdateThenPoll(ctx, id, cache); err != nil {
			return fmt.Errorf("Updating HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	// If any directory setting is set, we'll further check either the `usernameDownloaded` (for LDAP/Flat File), or the `domainJoined` (for AD) in response to ensure the configuration is correct, and the cache is functional.
	// There are situations that the LRO succeeded, whilst ends up with a non-functional cache (e.g. providing some invalid flat file setting).
	if directorySetting != nil {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving HPC Cache %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		model := resp.Model
		if model == nil {
			return fmt.Errorf("Unepxected nil `cacheProperties` in response")
		}
		prop := model.Properties
		if prop == nil {
			return fmt.Errorf("Unepxected nil `cacheProperties` in response")
		}
		ds := prop.DirectoryServicesSettings
		if ds == nil {
			return fmt.Errorf("Unexpected nil `directoryServicesSettings` in response")
		}

		// In case the user uses active directory service, we
		if directorySetting.ActiveDirectory != nil {
			ad := ds.ActiveDirectory
			if ad == nil || ad.DomainJoined == nil {
				return fmt.Errorf("Unexpected nil `activeDirectory` in response")
			}
			if *ad.DomainJoined != caches.DomainJoinedTypeYes {
				return fmt.Errorf("failed to join domain, current status: %s", *ad.DomainJoined)
			}
		} else {
			ud := ds.UsernameDownload
			if ud == nil || ud.UsernameDownloaded == nil {
				return fmt.Errorf("Unexpected nil `usernameDownload` in response")
			}
			if *ud.UsernameDownloaded != caches.UsernameDownloadedTypeYes {
				return fmt.Errorf("failed to download directory info, current status: %s", *ud.UsernameDownloaded)
			}
		}
	}

	d.SetId(id.ID())

	// wait for HPC Cache provision state to be succeeded. or further operations with it may fail.
	cacheClient := meta.(*clients.Client).StorageCache.Caches
	if _, err = resourceHPCCacheWaitForCreating(ctx, cacheClient, id, d); err != nil {
		return fmt.Errorf("waiting for the HPC Cache provision state %s (Resource Group: %s) : %+v", name, resourceGroup, err)
	}

	return resourceHPCCacheRead(d, meta)
}

func resourceHPCCacheRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := caches.ParseCacheID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] HPC Cache was not found in - removing from state! (%s)", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HPC Cache %q: %+v", id, err)
	}

	d.Set("name", id.CacheName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if m := resp.Model; m != nil {
		if sku := m.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		if props := m.Properties; props != nil {
			d.Set("location", azure.NormalizeLocation(pointer.From(m.Location)))
			d.Set("cache_size_in_gb", props.CacheSizeGB)
			d.Set("subnet_id", props.Subnet)
			d.Set("mount_addresses", utils.FlattenStringSlice(props.MountAddresses))

			mtu, ntpServer, dnsSetting := flattenStorageCacheNetworkSettings(props.NetworkSettings)
			d.Set("mtu", mtu)
			d.Set("ntp_server", ntpServer)
			if err := d.Set("dns", dnsSetting); err != nil {
				return fmt.Errorf("setting `dns`: %v", err)
			}

			ad, flatFile, ldap, err := flattenStorageCacheDirectorySettings(d, props.DirectoryServicesSettings)
			if err != nil {
				return err
			}

			if err := d.Set("directory_active_directory", ad); err != nil {
				return fmt.Errorf("setting `directory_active_directory`: %v", err)
			}

			if err := d.Set("directory_flat_file", flatFile); err != nil {
				return fmt.Errorf("setting `directory_flat_file`: %v", err)
			}

			if err := d.Set("directory_ldap", ldap); err != nil {
				return fmt.Errorf("setting `directory_ldap`: %v", err)
			}

			if securitySettings := props.SecuritySettings; securitySettings != nil {
				if securitySettings.AccessPolicies != nil {
					defaultPolicy := CacheGetAccessPolicyByName(*securitySettings.AccessPolicies, "default")
					if defaultPolicy != nil {
						defaultAccessPolicy, err := flattenStorageCacheNfsDefaultAccessPolicy(*defaultPolicy)
						if err != nil {
							return err
						}
						if err := d.Set("default_access_policy", defaultAccessPolicy); err != nil {
							return fmt.Errorf("setting `default_access_policy`: %v", err)
						}
					}
				}
			}

			keyVaultKeyId := ""
			autoKeyRotationEnabled := false
			if eprops := props.EncryptionSettings; eprops != nil {
				if eprops.KeyEncryptionKey != nil {
					keyVaultKeyId = eprops.KeyEncryptionKey.KeyURL
				}

				if eprops.RotationToLatestKeyVersionEnabled != nil {
					autoKeyRotationEnabled = *eprops.RotationToLatestKeyVersionEnabled
				}
			}
			d.Set("key_vault_key_id", keyVaultKeyId)
			d.Set("automatically_rotate_key_to_latest_enabled", autoKeyRotationEnabled)
		}

		i, err := identity.FlattenSystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, m.Tags)
	}

	return nil
}

func resourceHPCCacheDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StorageCache.Caches
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := caches.ParseCacheID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting HPC Cache %s: %+v", id, err)
	}

	return nil
}

func expandStorageCacheDefaultAccessPolicy(input []interface{}) *caches.NfsAccessPolicy {
	if len(input) == 0 {
		return nil
	}

	return &caches.NfsAccessPolicy{
		Name:        "default",
		AccessRules: expandStorageCacheNfsAccessRules(input[0].(map[string]interface{})["access_rule"].(*pluginsdk.Set).List()),
	}
}

func flattenStorageCacheNfsDefaultAccessPolicy(input caches.NfsAccessPolicy) ([]interface{}, error) {
	rules, err := flattenStorageCacheNfsAccessRules(input.AccessRules)
	if err != nil {
		return nil, err
	}

	return []interface{}{
		map[string]interface{}{
			"access_rule": rules,
		},
	}, nil
}

func expandStorageCacheNfsAccessRules(input []interface{}) []caches.NfsAccessRule {
	var out []caches.NfsAccessRule
	for _, accessRuleRaw := range input {
		b := accessRuleRaw.(map[string]interface{})
		out = append(out, caches.NfsAccessRule{
			Scope:          caches.NfsAccessRuleScope(b["scope"].(string)),
			Access:         caches.NfsAccessRuleAccess(b["access"].(string)),
			Filter:         pointer.To(b["filter"].(string)),
			Suid:           pointer.To(b["suid_enabled"].(bool)),
			SubmountAccess: pointer.To(b["submount_access_enabled"].(bool)),
			RootSquash:     pointer.To(b["root_squash_enabled"].(bool)),
			AnonymousUID:   pointer.To(strconv.Itoa(b["anonymous_uid"].(int))),
			AnonymousGID:   pointer.To(strconv.Itoa(b["anonymous_gid"].(int))),
		})
	}
	return out
}

func flattenStorageCacheNfsAccessRules(input []caches.NfsAccessRule) ([]interface{}, error) {
	var rules []interface{}
	for _, accessRule := range input {

		anonymousUID := 0
		if accessRule.AnonymousUID != nil {
			var err error
			anonymousUID, err = strconv.Atoi(*accessRule.AnonymousUID)
			if err != nil {
				return nil, fmt.Errorf("converting `anonymous_uid` from string to int")
			}
		}

		anonymousGID := 0
		if accessRule.AnonymousGID != nil {
			var err error
			anonymousGID, err = strconv.Atoi(*accessRule.AnonymousGID)
			if err != nil {
				return nil, fmt.Errorf("converting `anonymous_gid` from string to int")
			}
		}

		rules = append(rules, map[string]interface{}{
			"scope":                   accessRule.Scope,
			"access":                  accessRule.Access,
			"filter":                  pointer.From(accessRule.Filter),
			"suid_enabled":            pointer.From(accessRule.Suid),
			"submount_access_enabled": pointer.From(accessRule.SubmountAccess),
			"root_squash_enabled":     pointer.From(accessRule.RootSquash),
			"anonymous_uid":           anonymousUID,
			"anonymous_gid":           anonymousGID,
		})
	}

	return rules, nil
}

func expandStorageCacheNetworkSettings(d *pluginsdk.ResourceData) *caches.CacheNetworkSettings {
	out := &caches.CacheNetworkSettings{
		Mtu:       utils.Int64(int64(d.Get("mtu").(int))),
		NtpServer: pointer.To(d.Get("ntp_server").(string)),
	}

	if dnsSetting, ok := d.GetOk("dns"); ok {
		dnsSetting := dnsSetting.([]interface{})[0].(map[string]interface{})
		out.DnsServers = utils.ExpandStringSlice(dnsSetting["servers"].([]interface{}))
		searchDomain := dnsSetting["search_domain"].(string)
		if searchDomain != "" {
			out.DnsSearchDomain = &searchDomain
		}
	}
	return out
}

func flattenStorageCacheNetworkSettings(settings *caches.CacheNetworkSettings) (mtu int, ntpServer string, dnsSetting []interface{}) {
	if settings == nil {
		return
	}

	mtu = int(pointer.From(settings.Mtu))
	ntpServer = pointer.From(settings.NtpServer)

	if settings.DnsServers != nil {
		dnsSetting = []interface{}{
			map[string]interface{}{
				"servers":       utils.FlattenStringSlice(settings.DnsServers),
				"search_domain": pointer.From(settings.DnsSearchDomain),
			},
		}
	}
	return
}

func expandStorageCacheDirectorySettings(d *pluginsdk.ResourceData) *caches.CacheDirectorySettings {
	if raw := d.Get("directory_active_directory").([]interface{}); len(raw) != 0 {
		b := raw[0].(map[string]interface{})

		var secondaryDNSPtr *string
		if secondaryDNS := b["dns_secondary_ip"].(string); secondaryDNS != "" {
			secondaryDNSPtr = &secondaryDNS
		}

		return &caches.CacheDirectorySettings{
			UsernameDownload: &caches.CacheUsernameDownloadSettings{
				ExtendedGroups: pointer.To(true),
				UsernameSource: pointer.To(caches.UsernameSourceAD),
			},
			ActiveDirectory: &caches.CacheActiveDirectorySettings{
				PrimaryDnsIPAddress:   b["dns_primary_ip"].(string),
				SecondaryDnsIPAddress: secondaryDNSPtr,
				DomainName:            b["domain_name"].(string),
				CacheNetBiosName:      b["cache_netbios_name"].(string),
				DomainNetBiosName:     b["domain_netbios_name"].(string),
				Credentials: &caches.CacheActiveDirectorySettingsCredentials{
					Username: b["username"].(string),
					Password: pointer.To(b["password"].(string)),
				},
			},
		}
	}

	if raw := d.Get("directory_flat_file").([]interface{}); len(raw) != 0 {
		b := raw[0].(map[string]interface{})
		return &caches.CacheDirectorySettings{
			UsernameDownload: &caches.CacheUsernameDownloadSettings{
				ExtendedGroups: pointer.To(true),
				UsernameSource: pointer.To(caches.UsernameSourceFile),
				GroupFileURI:   pointer.To(b["group_file_uri"].(string)),
				UserFileURI:    pointer.To(b["password_file_uri"].(string)),
			},
		}
	}

	if raw := d.Get("directory_ldap").([]interface{}); len(raw) != 0 {
		b := raw[0].(map[string]interface{})
		var certValidationUriPtr *string
		certValidationUri := b["certificate_validation_uri"].(string)
		if certValidationUri != "" {
			certValidationUriPtr = &certValidationUri
		}
		return &caches.CacheDirectorySettings{
			UsernameDownload: &caches.CacheUsernameDownloadSettings{
				ExtendedGroups:          pointer.To(true),
				UsernameSource:          pointer.To(caches.UsernameSourceLDAP),
				LdapServer:              pointer.To(b["server"].(string)),
				LdapBaseDN:              pointer.To(b["base_dn"].(string)),
				EncryptLdapConnection:   pointer.To(b["encrypted"].(bool)),
				RequireValidCertificate: pointer.To(certValidationUriPtr != nil),
				AutoDownloadCertificate: pointer.To(b["download_certificate_automatically"].(bool)),
				CaCertificateURI:        certValidationUriPtr,
				Credentials:             expandStorageCacheDirectoryLdapBind(b["bind"].([]interface{})),
			},
		}
	}

	return nil
}

func flattenStorageCacheDirectorySettings(d *pluginsdk.ResourceData, input *caches.CacheDirectorySettings) (ad, flatFile, ldap []interface{}, err error) {
	if input == nil || input.UsernameDownload == nil || *input.UsernameDownload.UsernameSource == caches.UsernameSourceNone {
		return nil, nil, nil, nil
	}

	ud := input.UsernameDownload
	switch *ud.UsernameSource {
	case caches.UsernameSourceAD:
		var (
			primaryDNS        string
			domainName        string
			cacheNetBiosName  string
			username          string
			password          string
			domainNetBiosName string
			secondaryDNS      string
		)

		if ad := input.ActiveDirectory; ad != nil {
			primaryDNS = ad.PrimaryDnsIPAddress
			domainName = ad.DomainName
			cacheNetBiosName = ad.CacheNetBiosName
			domainNetBiosName = ad.DomainNetBiosName
			secondaryDNS = pointer.From(ad.SecondaryDnsIPAddress)
		}
		// Since the credentials are never returned from response. We will set whatever specified in the config back to state as the best effort.
		ad := d.Get("directory_active_directory").([]interface{})
		if len(ad) == 1 {
			b := ad[0].(map[string]interface{})
			username = b["username"].(string)
			password = b["password"].(string)
		}

		return []interface{}{
			map[string]interface{}{
				"dns_primary_ip":      primaryDNS,
				"domain_name":         domainName,
				"cache_netbios_name":  cacheNetBiosName,
				"domain_netbios_name": domainNetBiosName,
				"dns_secondary_ip":    secondaryDNS,
				"username":            username,
				"password":            password,
			},
		}, nil, nil, nil

	case caches.UsernameSourceFile:
		return nil, []interface{}{
			map[string]interface{}{
				"group_file_uri":    pointer.From(ud.GroupFileURI),
				"password_file_uri": pointer.From(ud.UserFileURI),
			},
		}, nil, nil
	case caches.UsernameSourceLDAP:
		return nil, nil, []interface{}{
			map[string]interface{}{
				"server":                             pointer.From(ud.LdapServer),
				"base_dn":                            pointer.From(ud.LdapBaseDN),
				"encrypted":                          pointer.From(ud.EncryptLdapConnection),
				"certificate_validation_uri":         pointer.From(ud.CaCertificateURI),
				"download_certificate_automatically": pointer.From(ud.AutoDownloadCertificate),
				"bind":                               flattenStorageCacheDirectoryLdapBind(d),
			},
		}, nil
	default:
		return nil, nil, nil, fmt.Errorf("source type %q is not supported", *ud.UsernameSource)
	}
}

func flattenStorageCacheDirectoryLdapBind(d *pluginsdk.ResourceData) []interface{} {
	// Since the credentials are never returned from response. We will set whatever specified in the config back to state as the best effort.
	ldap := d.Get("directory_ldap").([]interface{})
	if len(ldap) == 0 {
		return nil
	}

	return ldap[0].(map[string]interface{})["bind"].([]interface{})
}

func expandStorageCacheDirectoryLdapBind(input []interface{}) *caches.CacheUsernameDownloadSettingsCredentials {
	if len(input) == 0 {
		return nil
	}

	b := input[0].(map[string]interface{})
	return &caches.CacheUsernameDownloadSettingsCredentials{
		BindDn:       pointer.To(b["dn"].(string)),
		BindPassword: pointer.To(b["password"].(string)),
	}
}

type storageCacheKeyVault struct {
	keyVaultId             string
	resourceGroupName      string
	keyVaultName           string
	location               string
	purgeProtectionEnabled bool
	softDeleteEnabled      bool
}

func storageCacheRetrieveKeyVault(ctx context.Context, keyVaultsClient *client.Client, subscriptionId string, id string) (*storageCacheKeyVault, error) {
	keyVaultKeyId, err := keyVaultParse.ParseNestedItemID(id)
	if err != nil {
		return nil, err
	}
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultKeyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseUrl, err)
	}
	if keyVaultID == nil {
		return nil, fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", keyVaultKeyId.KeyVaultBaseUrl)
	}

	parsedKeyVaultID, err := commonids.ParseKeyVaultID(*keyVaultID)
	if err != nil {
		return nil, err
	}

	resp, err := keyVaultsClient.VaultsClient.Get(ctx, *parsedKeyVaultID)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *parsedKeyVaultID, err)
	}

	loc := ""
	purgeProtectionEnabled := false
	softDeleteEnabled := false
	if model := resp.Model; model != nil {
		loc = location.NormalizeNilable(model.Location)
		if model.Properties.EnableSoftDelete != nil {
			softDeleteEnabled = *model.Properties.EnableSoftDelete
		}

		if model.Properties.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *model.Properties.EnablePurgeProtection
		}
	}

	return &storageCacheKeyVault{
		keyVaultId:             *keyVaultID,
		resourceGroupName:      parsedKeyVaultID.ResourceGroupName,
		keyVaultName:           parsedKeyVaultID.VaultName,
		location:               loc,
		purgeProtectionEnabled: purgeProtectionEnabled,
		softDeleteEnabled:      softDeleteEnabled,
	}, nil
}

func resourceHPCCacheSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cache_size_in_gb": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.IntInSlice([]int{
				3072,
				6144,
				12288,
				21623,
				24576,
				43246,
				49152,
				86491,
			}),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard_2G",
				"Standard_4G",
				"Standard_8G",
				"Standard_L4_5G",
				"Standard_L9G",
				"Standard_L16G",
			}, false),
		},

		"mtu": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1500,
			ValidateFunc: validation.IntBetween(576, 1500),
		},

		"ntp_server": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "time.windows.com",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"servers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 3,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPAddress,
						},
					},

					"search_domain": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"directory_active_directory": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_primary_ip": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPAddress,
					},
					"domain_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"cache_netbios_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[-0-9a-zA-Z]{1,15}$`), ""),
					},
					"domain_netbios_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[-0-9a-zA-Z]{1,15}$`), ""),
					},
					"username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"dns_secondary_ip": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPAddress,
					},
				},
			},
			ConflictsWith: []string{"directory_flat_file", "directory_ldap"},
		},

		"directory_flat_file": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_file_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password_file_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ConflictsWith: []string{"directory_active_directory", "directory_ldap"},
		},

		"directory_ldap": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"base_dn": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"encrypted": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"certificate_validation_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"download_certificate_automatically": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						RequiredWith: []string{"directory_ldap.0.certificate_validation_uri"},
					},

					"bind": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"password": {
									Type:         pluginsdk.TypeString,
									Sensitive:    true,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
			ConflictsWith: []string{"directory_active_directory", "directory_flat_file"},
		},

		"default_access_policy": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Optional: true,
			// This is computed because there is always a "default" policy in the cache. It is created together with the cache, and users can't remove it.
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"access_rule": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						MaxItems: 3,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"scope": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(caches.NfsAccessRuleScopeDefault),
										string(caches.NfsAccessRuleScopeNetwork),
										string(caches.NfsAccessRuleScopeHost),
									}, false),
								},

								"access": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(caches.NfsAccessRuleAccessRw),
										string(caches.NfsAccessRuleAccessRo),
										string(caches.NfsAccessRuleAccessNo),
									}, false),
								},

								"filter": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"suid_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"submount_access_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"root_squash_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"anonymous_uid": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"anonymous_gid": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},
							},
						},
					},
				},
			},
		},

		"mount_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.NestedItemId,
			RequiredWith: []string{"identity"},
		},

		"automatically_rotate_key_to_latest_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"key_vault_key_id"},
		},

		"tags": commonschema.Tags(),
	}
}
