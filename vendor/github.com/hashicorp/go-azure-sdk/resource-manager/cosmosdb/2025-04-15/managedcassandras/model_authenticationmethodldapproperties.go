package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMethodLdapProperties struct {
	ConnectionTimeoutInMs        *int64         `json:"connectionTimeoutInMs,omitempty"`
	SearchBaseDistinguishedName  *string        `json:"searchBaseDistinguishedName,omitempty"`
	SearchFilterTemplate         *string        `json:"searchFilterTemplate,omitempty"`
	ServerCertificates           *[]Certificate `json:"serverCertificates,omitempty"`
	ServerHostname               *string        `json:"serverHostname,omitempty"`
	ServerPort                   *int64         `json:"serverPort,omitempty"`
	ServiceUserDistinguishedName *string        `json:"serviceUserDistinguishedName,omitempty"`
	ServiceUserPassword          *string        `json:"serviceUserPassword,omitempty"`
}
