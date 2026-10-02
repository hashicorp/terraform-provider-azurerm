package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureActiveDirectoryRegistration struct {
	ClientId                                      *string `json:"clientId,omitempty"`
	ClientSecretCertificateIssuer                 *string `json:"clientSecretCertificateIssuer,omitempty"`
	ClientSecretCertificateSubjectAlternativeName *string `json:"clientSecretCertificateSubjectAlternativeName,omitempty"`
	ClientSecretCertificateThumbprint             *string `json:"clientSecretCertificateThumbprint,omitempty"`
	ClientSecretSettingName                       *string `json:"clientSecretSettingName,omitempty"`
	OpenIdIssuer                                  *string `json:"openIdIssuer,omitempty"`
}
