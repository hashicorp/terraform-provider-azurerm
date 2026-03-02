package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeConnectionInfo struct {
	HostServiceUri         *string `json:"hostServiceUri,omitempty"`
	IdentityCertThumbprint *string `json:"identityCertThumbprint,omitempty"`
	IsIdentityCertExprired *bool   `json:"isIdentityCertExprired,omitempty"`
	PublicKey              *string `json:"publicKey,omitempty"`
	ServiceToken           *string `json:"serviceToken,omitempty"`
	Version                *string `json:"version,omitempty"`
}
