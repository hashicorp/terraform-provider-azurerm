package vmmservers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMmServerProperties struct {
	ConnectionStatus  *string            `json:"connectionStatus,omitempty"`
	Credentials       *VMmCredential     `json:"credentials,omitempty"`
	ErrorMessage      *string            `json:"errorMessage,omitempty"`
	Fqdn              string             `json:"fqdn"`
	Port              *int64             `json:"port,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Uuid              *string            `json:"uuid,omitempty"`
	Version           *string            `json:"version,omitempty"`
}
