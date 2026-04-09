package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminCredentials struct {
	NsxtPassword    *string `json:"nsxtPassword,omitempty"`
	NsxtUsername    *string `json:"nsxtUsername,omitempty"`
	VcenterPassword *string `json:"vcenterPassword,omitempty"`
	VcenterUsername *string `json:"vcenterUsername,omitempty"`
}
