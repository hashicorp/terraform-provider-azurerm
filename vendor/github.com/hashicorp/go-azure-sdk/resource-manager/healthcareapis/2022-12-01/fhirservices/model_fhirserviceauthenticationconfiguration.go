package fhirservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FhirServiceAuthenticationConfiguration struct {
	Audience          *string `json:"audience,omitempty"`
	Authority         *string `json:"authority,omitempty"`
	SmartProxyEnabled *bool   `json:"smartProxyEnabled,omitempty"`
}
