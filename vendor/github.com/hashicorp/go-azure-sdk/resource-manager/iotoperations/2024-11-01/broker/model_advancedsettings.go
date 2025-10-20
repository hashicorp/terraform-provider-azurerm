package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedSettings struct {
	Clients                *ClientConfig           `json:"clients,omitempty"`
	EncryptInternalTraffic *OperationalMode        `json:"encryptInternalTraffic,omitempty"`
	InternalCerts          *CertManagerCertOptions `json:"internalCerts,omitempty"`
}
