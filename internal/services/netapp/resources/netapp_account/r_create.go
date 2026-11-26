// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_account

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.AccountClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := netappaccounts.NewNetAppAccountID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_account", id.ID())
		}
	}

	accountParameters := netappaccounts.NetAppAccount{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &netappaccounts.AccountProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	activeDirectoryRaw := d.Get("active_directory")
	if activeDirectoryRaw != nil {
		activeDirectories := activeDirectoryRaw.([]any)
		activeDirectoriesExpanded := expandNetAppActiveDirectories(activeDirectories)
		if len(pointer.From(activeDirectoriesExpanded)) > 0 {
			accountParameters.Properties.ActiveDirectories = activeDirectoriesExpanded
		}
	}

	anfAccountIdentityRaw := d.Get("identity")
	if anfAccountIdentityRaw != nil {
		anfAccountIdentity, ok := anfAccountIdentityRaw.([]any)

		if ok && len(anfAccountIdentity) > 0 {
			anfAccountIdentityExpanded, err := identity.ExpandLegacySystemAndUserAssignedMap(anfAccountIdentity)
			if err != nil {
				return err
			}
			if anfAccountIdentity != nil {
				accountParameters.Identity = anfAccountIdentityExpanded
			}
		}
	}

	if err := client.AccountsCreateOrUpdateCallbackThenPoll(ctx, id, accountParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetAppAccountRead(d, meta)
}

func expandNetAppActiveDirectories(input []any) *[]netappaccounts.ActiveDirectory {
	results := make([]netappaccounts.ActiveDirectory, 0)
	if input == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]any)
		dns := strings.Join(*pluginsdk.ExpandStringSlice(v["dns_servers"].([]any)), ",")

		result := netappaccounts.ActiveDirectory{
			Dns:                           pointer.To(dns),
			Domain:                        pointer.To(v["domain"].(string)),
			OrganizationalUnit:            pointer.To(v["organizational_unit"].(string)),
			Password:                      pointer.To(v["password"].(string)),
			SmbServerName:                 pointer.To(v["smb_server_name"].(string)),
			Username:                      pointer.To(v["username"].(string)),
			Site:                          pointer.To(v["site_name"].(string)),
			AdName:                        pointer.To(v["kerberos_ad_name"].(string)),
			KdcIP:                         pointer.To(v["kerberos_kdc_ip"].(string)),
			AesEncryption:                 pointer.To(v["aes_encryption_enabled"].(bool)),
			AllowLocalNfsUsersWithLdap:    pointer.To(v["local_nfs_users_with_ldap_allowed"].(bool)),
			LdapOverTLS:                   pointer.To(v["ldap_over_tls_enabled"].(bool)),
			ServerRootCACertificate:       pointer.To(v["server_root_ca_certificate"].(string)),
			LdapSigning:                   pointer.To(v["ldap_signing_enabled"].(bool)),
			EncryptDCConnections:          pointer.To(v["domain_controller_connection_encryption_enabled"].(bool)),
			PreferredServersForLdapClient: pointer.To(strings.Join(*pluginsdk.ExpandStringSlice(v["ldap_client_preferred_servers"].([]any)), ",")),
		}

		var searchScope map[string]any
		if rawSearchScope := v["ldap_search_scope"].([]any); len(rawSearchScope) > 0 {
			searchScope = rawSearchScope[0].(map[string]any)
		}
		result.LdapSearchScope = expandNetAppLdapSearchScope(searchScope)

		if admins := v["administrator_users"].(*pluginsdk.Set); admins != nil {
			result.Administrators = pluginsdk.ExpandStringSlice(admins.List())
		}

		if backupUsers := v["backup_users"].(*pluginsdk.Set); backupUsers != nil {
			result.BackupOperators = pluginsdk.ExpandStringSlice(backupUsers.List())
		}

		if secUsers := v["security_users"].(*pluginsdk.Set); secUsers != nil {
			result.SecurityOperators = pluginsdk.ExpandStringSlice(secUsers.List())
		}

		results = append(results, result)
	}
	return &results
}

func expandNetAppLdapSearchScope(input map[string]any) *netappaccounts.LdapSearchScopeOpt {
	if input == nil {
		// if not supplied, set empty strings to all fields to properly update in API.
		// if LdapSearchScopeOpt is left nil or empty, the API will not update the search scope and will leave it as-is.
		return &netappaccounts.LdapSearchScopeOpt{
			GroupDN:               pointer.To(""),
			GroupMembershipFilter: pointer.To(""),
			UserDN:                pointer.To(""),
		}
	}

	result := netappaccounts.LdapSearchScopeOpt{
		GroupDN:               pointer.To(input["group_dn"].(string)),
		GroupMembershipFilter: pointer.To(input["group_membership_filter"].(string)),
		UserDN:                pointer.To(input["user_dn"].(string)),
	}

	return &result
}
