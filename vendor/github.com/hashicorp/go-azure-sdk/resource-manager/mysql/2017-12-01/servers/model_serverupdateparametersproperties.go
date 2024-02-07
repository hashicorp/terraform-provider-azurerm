package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerUpdateParametersProperties struct {
	AdministratorLoginPassword *string                  `json:"administratorLoginPassword,omitempty"`
	MinimalTlsVersion          *MinimalTlsVersionEnum   `json:"minimalTlsVersion,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccessEnum `json:"publicNetworkAccess,omitempty"`
	ReplicationRole            *string                  `json:"replicationRole,omitempty"`
	SslEnforcement             *SslEnforcementEnum      `json:"sslEnforcement,omitempty"`
	StorageProfile             *StorageProfile          `json:"storageProfile,omitempty"`
	Version                    *ServerVersion           `json:"version,omitempty"`
}
