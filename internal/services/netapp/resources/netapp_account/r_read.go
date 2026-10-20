// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		anfAccountIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", anfAccountIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if model.Properties != nil {
			// the API doesn't return values for `active_directory.0.password` and `active_directory.0.server_root_ca_certificate`, so we pass through current state values so change detection works
			prevPassword := d.Get("active_directory.0.password").(string)
			prevCaCert := d.Get("active_directory.0.server_root_ca_certificate").(string)

			if err = d.Set("active_directory", flattenNetAppActiveDirectories(model.Properties.ActiveDirectories, &prevPassword, &prevCaCert)); err != nil {
				return fmt.Errorf("setting `active_directory`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenNetAppActiveDirectories(input *[]netappaccounts.ActiveDirectory, prevPassword *string, prevCaCert *string) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}

	v := (*input)[0]

	return []any{
		map[string]any{
			"dns_servers":                       pluginsdk.FlattenStringSliceWithDelimiter(v.Dns, ","),
			"domain":                            v.Domain,
			"organizational_unit":               v.OrganizationalUnit,
			"password":                          prevPassword,
			"smb_server_name":                   v.SmbServerName,
			"username":                          v.Username,
			"site_name":                         v.Site,
			"kerberos_ad_name":                  v.AdName,
			"kerberos_kdc_ip":                   v.KdcIP,
			"aes_encryption_enabled":            v.AesEncryption,
			"local_nfs_users_with_ldap_allowed": v.AllowLocalNfsUsersWithLdap,
			"ldap_over_tls_enabled":             v.LdapOverTLS,
			"server_root_ca_certificate":        prevCaCert,
			"ldap_signing_enabled":              v.LdapSigning,
		},
	}
}
