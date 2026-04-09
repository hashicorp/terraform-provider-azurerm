package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityEndpoint struct {
	Location         *string `json:"location,omitempty"`
	Name             *string `json:"name,omitempty"`
	Port             *int64  `json:"port,omitempty"`
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
	Protocol         *string `json:"protocol,omitempty"`
}
