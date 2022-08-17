package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainSecuritySettings struct {
	KerberosArmoring      *KerberosArmoring      `json:"kerberosArmoring,omitempty"`
	KerberosRc4Encryption *KerberosRc4Encryption `json:"kerberosRc4Encryption,omitempty"`
	NtlmV1                *NtlmV1                `json:"ntlmV1,omitempty"`
	SyncKerberosPasswords *SyncKerberosPasswords `json:"syncKerberosPasswords,omitempty"`
	SyncNtlmPasswords     *SyncNtlmPasswords     `json:"syncNtlmPasswords,omitempty"`
	SyncOnPremPasswords   *SyncOnPremPasswords   `json:"syncOnPremPasswords,omitempty"`
	TlsV1                 *TlsV1                 `json:"tlsV1,omitempty"`
}
