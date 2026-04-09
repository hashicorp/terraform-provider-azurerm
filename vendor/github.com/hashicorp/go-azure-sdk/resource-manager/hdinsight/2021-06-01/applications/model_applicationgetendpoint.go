package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGetEndpoint struct {
	DestinationPort  *int64  `json:"destinationPort,omitempty"`
	Location         *string `json:"location,omitempty"`
	PrivateIPAddress *string `json:"privateIPAddress,omitempty"`
	PublicPort       *int64  `json:"publicPort,omitempty"`
}
