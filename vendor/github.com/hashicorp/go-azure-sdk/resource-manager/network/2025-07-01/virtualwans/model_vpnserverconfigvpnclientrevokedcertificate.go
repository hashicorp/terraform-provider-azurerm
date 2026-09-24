package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnServerConfigVpnClientRevokedCertificate struct {
	Name       *string `json:"name,omitempty"`
	Thumbprint *string `json:"thumbprint,omitempty"`
}
