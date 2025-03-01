package assetendpointprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Authentication struct {
	Method                      AuthenticationMethod         `json:"method"`
	UsernamePasswordCredentials *UsernamePasswordCredentials `json:"usernamePasswordCredentials,omitempty"`
	X509Credentials             *X509Credentials             `json:"x509Credentials,omitempty"`
}
