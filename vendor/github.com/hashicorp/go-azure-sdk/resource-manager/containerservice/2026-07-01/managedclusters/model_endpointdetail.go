package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointDetail struct {
	Description *string `json:"description,omitempty"`
	IPAddress   *string `json:"ipAddress,omitempty"`
	Port        *int64  `json:"port,omitempty"`
	Protocol    *string `json:"protocol,omitempty"`
}
