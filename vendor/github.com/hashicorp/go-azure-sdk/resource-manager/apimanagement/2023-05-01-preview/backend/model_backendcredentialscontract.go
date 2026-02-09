package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendCredentialsContract struct {
	Authorization  *BackendAuthorizationHeaderCredentials `json:"authorization,omitempty"`
	Certificate    *[]string                              `json:"certificate,omitempty"`
	CertificateIds *[]string                              `json:"certificateIds,omitempty"`
	Header         *map[string][]string                   `json:"header,omitempty"`
	Query          *map[string][]string                   `json:"query,omitempty"`
}
