package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheUsernameDownloadSettings struct {
	AutoDownloadCertificate *bool                                     `json:"autoDownloadCertificate,omitempty"`
	CaCertificateURI        *string                                   `json:"caCertificateURI,omitempty"`
	Credentials             *CacheUsernameDownloadSettingsCredentials `json:"credentials,omitempty"`
	EncryptLdapConnection   *bool                                     `json:"encryptLdapConnection,omitempty"`
	ExtendedGroups          *bool                                     `json:"extendedGroups,omitempty"`
	GroupFileURI            *string                                   `json:"groupFileURI,omitempty"`
	LdapBaseDN              *string                                   `json:"ldapBaseDN,omitempty"`
	LdapServer              *string                                   `json:"ldapServer,omitempty"`
	RequireValidCertificate *bool                                     `json:"requireValidCertificate,omitempty"`
	UserFileURI             *string                                   `json:"userFileURI,omitempty"`
	UsernameDownloaded      *UsernameDownloadedType                   `json:"usernameDownloaded,omitempty"`
	UsernameSource          *UsernameSource                           `json:"usernameSource,omitempty"`
}
