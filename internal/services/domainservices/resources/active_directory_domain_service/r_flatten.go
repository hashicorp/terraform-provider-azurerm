// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenDomainServiceLdaps(d *pluginsdk.ResourceData, input *domainservices.LdapsSettings, dataSource bool) []any {
	result := map[string]any{
		"enabled":                 false,
		"external_access_enabled": false,
		"certificate_expiry":      "",
		"certificate_thumbprint":  "",
		"public_certificate":      "",
	}

	if !dataSource {
		// Read pfx_certificate and pfx_certificate_password from existing state since it's not returned
		result["pfx_certificate"] = ""
		if v, ok := d.GetOk("secure_ldap.0.pfx_certificate"); ok {
			result["pfx_certificate"] = v.(string)
		}
		result["pfx_certificate_password"] = ""
		if v, ok := d.GetOk("secure_ldap.0.pfx_certificate_password"); ok {
			result["pfx_certificate_password"] = v.(string)
		}
	}

	if input != nil {
		if input.ExternalAccess != nil && *input.ExternalAccess == domainservices.ExternalAccessEnabled {
			result["external_access_enabled"] = true
		}
		if input.Ldaps != nil && *input.Ldaps == domainservices.LdapsEnabled {
			result["enabled"] = true
		}
		if v := input.CertificateNotAfter; v != nil {
			result["certificate_expiry"] = *v
		}
		if v := input.CertificateThumbprint; v != nil {
			result["certificate_thumbprint"] = *v
		}
		if v := input.PublicCertificate; v != nil {
			result["public_certificate"] = *v
		}
	}

	return []any{result}
}

func flattenDomainServiceNotifications(input *domainservices.NotificationSettings) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := map[string]any{
		"additional_recipients": make([]string, 0),
		"notify_dc_admins":      false,
		"notify_global_admins":  false,
	}
	if input.AdditionalRecipients != nil {
		result["additional_recipients"] = *input.AdditionalRecipients
	}
	if input.NotifyDcAdmins != nil && *input.NotifyDcAdmins == domainservices.NotifyDcAdminsEnabled {
		result["notify_dc_admins"] = true
	}
	if input.NotifyGlobalAdmins != nil && *input.NotifyGlobalAdmins == domainservices.NotifyGlobalAdminsEnabled {
		result["notify_global_admins"] = true
	}

	return []any{result}
}

func flattenDomainServiceSecurity(input *domainservices.DomainSecuritySettings) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := map[string]bool{
		"kerberos_armoring_enabled":       false,
		"kerberos_rc4_encryption_enabled": false,
		"ntlm_v1_enabled":                 false,
		"sync_kerberos_passwords":         false,
		"sync_ntlm_passwords":             false,
		"sync_on_prem_passwords":          false,
		"tls_v1_enabled":                  false,
	}
	if input.KerberosArmoring != nil && *input.KerberosArmoring == domainservices.KerberosArmoringEnabled {
		result["kerberos_armoring_enabled"] = true
	}
	if input.KerberosRc4Encryption != nil && *input.KerberosRc4Encryption == domainservices.KerberosRc4EncryptionEnabled {
		result["kerberos_rc4_encryption_enabled"] = true
	}
	if input.NtlmV1 != nil && *input.NtlmV1 == domainservices.NtlmV1Enabled {
		result["ntlm_v1_enabled"] = true
	}
	if input.SyncKerberosPasswords != nil && *input.SyncKerberosPasswords == domainservices.SyncKerberosPasswordsEnabled {
		result["sync_kerberos_passwords"] = true
	}
	if input.SyncNtlmPasswords != nil && *input.SyncNtlmPasswords == domainservices.SyncNtlmPasswordsEnabled {
		result["sync_ntlm_passwords"] = true
	}
	if input.SyncOnPremPasswords != nil && *input.SyncOnPremPasswords == domainservices.SyncOnPremPasswordsEnabled {
		result["sync_on_prem_passwords"] = true
	}
	if input.TlsV1 != nil && *input.TlsV1 == domainservices.TlsV1Enabled {
		result["tls_v1_enabled"] = true
	}

	return []any{result}
}
