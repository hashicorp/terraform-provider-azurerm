package attestationproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StatusResult struct {
	AttestUri                  *string                      `json:"attestUri,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	Status                     *AttestationServiceStatus    `json:"status,omitempty"`
	TrustModel                 *string                      `json:"trustModel,omitempty"`
}
