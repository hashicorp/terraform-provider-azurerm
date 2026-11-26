// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandDomainServiceLdaps(input []any) (ldaps *domainservices.LdapsSettings) {
	ldaps = &domainservices.LdapsSettings{
		Ldaps: pointer.To(domainservices.LdapsDisabled),
	}

	if len(input) > 0 {
		v := input[0].(map[string]any)
		if v["enabled"].(bool) {
			*ldaps.Ldaps = domainservices.LdapsEnabled
		}
		ldaps.PfxCertificate = pointer.To(v["pfx_certificate"].(string))
		ldaps.PfxCertificatePassword = pointer.To(v["pfx_certificate_password"].(string))
		access := domainservices.ExternalAccessDisabled
		if v["external_access_enabled"].(bool) {
			access = domainservices.ExternalAccessEnabled
		}
		ldaps.ExternalAccess = &access
	}

	return
}

func expandDomainServiceNotifications(input []any) *domainservices.NotificationSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	additionalRecipients := make([]string, 0)
	if ar, ok := v["additional_recipients"]; ok {
		for _, r := range ar.(*pluginsdk.Set).List() {
			additionalRecipients = append(additionalRecipients, r.(string))
		}
	}

	notifyDcAdmins := domainservices.NotifyDcAdminsDisabled
	if n, ok := v["notify_dc_admins"]; ok && n.(bool) {
		notifyDcAdmins = domainservices.NotifyDcAdminsEnabled
	}

	notifyGlobalAdmins := domainservices.NotifyGlobalAdminsDisabled
	if n, ok := v["notify_global_admins"]; ok && n.(bool) {
		notifyGlobalAdmins = domainservices.NotifyGlobalAdminsEnabled
	}

	return &domainservices.NotificationSettings{
		AdditionalRecipients: &additionalRecipients,
		NotifyDcAdmins:       &notifyDcAdmins,
		NotifyGlobalAdmins:   &notifyGlobalAdmins,
	}
}

func expandDomainServiceSecurity(input []any) *domainservices.DomainSecuritySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	kerberosRc4Encryption := domainservices.KerberosRc4EncryptionDisabled
	kerberosArmoring := domainservices.KerberosArmoringDisabled
	ntlmV1 := domainservices.NtlmV1Disabled
	syncKerberosPasswords := domainservices.SyncKerberosPasswordsDisabled
	syncNtlmPasswords := domainservices.SyncNtlmPasswordsDisabled
	syncOnPremPasswords := domainservices.SyncOnPremPasswordsDisabled
	tlsV1 := domainservices.TlsV1Disabled

	if v["kerberos_armoring_enabled"].(bool) {
		kerberosArmoring = domainservices.KerberosArmoringEnabled
	}
	if v["kerberos_rc4_encryption_enabled"].(bool) {
		kerberosRc4Encryption = domainservices.KerberosRc4EncryptionEnabled
	}
	if v["ntlm_v1_enabled"].(bool) {
		ntlmV1 = domainservices.NtlmV1Enabled
	}
	if v["sync_kerberos_passwords"].(bool) {
		syncKerberosPasswords = domainservices.SyncKerberosPasswordsEnabled
	}
	if v["sync_ntlm_passwords"].(bool) {
		syncNtlmPasswords = domainservices.SyncNtlmPasswordsEnabled
	}
	if v["sync_on_prem_passwords"].(bool) {
		syncOnPremPasswords = domainservices.SyncOnPremPasswordsEnabled
	}
	if v["tls_v1_enabled"].(bool) {
		tlsV1 = domainservices.TlsV1Enabled
	}

	return &domainservices.DomainSecuritySettings{
		KerberosArmoring:      &kerberosArmoring,
		KerberosRc4Encryption: &kerberosRc4Encryption,
		NtlmV1:                &ntlmV1,
		SyncKerberosPasswords: &syncKerberosPasswords,
		SyncNtlmPasswords:     &syncNtlmPasswords,
		SyncOnPremPasswords:   &syncOnPremPasswords,
		TlsV1:                 &tlsV1,
	}
}
