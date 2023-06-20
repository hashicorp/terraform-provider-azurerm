package managedidentities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FederatedIdentityCredentialProperties struct {
	Audiences []string `json:"audiences"`
	Issuer    string   `json:"issuer"`
	Subject   string   `json:"subject"`
}
