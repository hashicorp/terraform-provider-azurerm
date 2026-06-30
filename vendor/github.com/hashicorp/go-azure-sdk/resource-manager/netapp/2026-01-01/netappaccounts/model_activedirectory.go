package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveDirectory struct {
	ActiveDirectoryId             *string                `json:"activeDirectoryId,omitempty"`
	AdName                        *string                `json:"adName,omitempty"`
	Administrators                *[]string              `json:"administrators,omitempty"`
	AesEncryption                 *bool                  `json:"aesEncryption,omitempty"`
	AllowLocalNfsUsersWithLdap    *bool                  `json:"allowLocalNfsUsersWithLdap,omitempty"`
	BackupOperators               *[]string              `json:"backupOperators,omitempty"`
	Dns                           *string                `json:"dns,omitempty"`
	Domain                        *string                `json:"domain,omitempty"`
	EncryptDCConnections          *bool                  `json:"encryptDCConnections,omitempty"`
	KdcIP                         *string                `json:"kdcIP,omitempty"`
	LdapOverTLS                   *bool                  `json:"ldapOverTLS,omitempty"`
	LdapSearchScope               *LdapSearchScopeOpt    `json:"ldapSearchScope,omitempty"`
	LdapSigning                   *bool                  `json:"ldapSigning,omitempty"`
	OrganizationalUnit            *string                `json:"organizationalUnit,omitempty"`
	Password                      *string                `json:"password,omitempty"`
	PreferredServersForLdapClient *string                `json:"preferredServersForLdapClient,omitempty"`
	SecurityOperators             *[]string              `json:"securityOperators,omitempty"`
	ServerRootCACertificate       *string                `json:"serverRootCACertificate,omitempty"`
	Site                          *string                `json:"site,omitempty"`
	SmbServerName                 *string                `json:"smbServerName,omitempty"`
	Status                        *ActiveDirectoryStatus `json:"status,omitempty"`
	StatusDetails                 *string                `json:"statusDetails,omitempty"`
	Username                      *string                `json:"username,omitempty"`
}
