package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssuerParameters struct {
	CertTransparency *bool   `json:"cert_transparency,omitempty"`
	Cty              *string `json:"cty,omitempty"`
	Name             *string `json:"name,omitempty"`
}
