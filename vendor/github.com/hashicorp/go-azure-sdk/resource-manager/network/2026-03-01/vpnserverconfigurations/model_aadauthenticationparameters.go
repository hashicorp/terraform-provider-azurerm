package vpnserverconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AadAuthenticationParameters struct {
	AadAudience *string `json:"aadAudience,omitempty"`
	AadIssuer   *string `json:"aadIssuer,omitempty"`
	AadTenant   *string `json:"aadTenant,omitempty"`
}
