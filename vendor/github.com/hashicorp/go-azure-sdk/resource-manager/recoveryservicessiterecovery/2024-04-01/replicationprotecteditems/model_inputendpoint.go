package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputEndpoint struct {
	EndpointName *string `json:"endpointName,omitempty"`
	PrivatePort  *int64  `json:"privatePort,omitempty"`
	Protocol     *string `json:"protocol,omitempty"`
	PublicPort   *int64  `json:"publicPort,omitempty"`
}
