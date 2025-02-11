package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmFailbackNicDetails struct {
	AdapterType     *string `json:"adapterType,omitempty"`
	MacAddress      *string `json:"macAddress,omitempty"`
	NetworkName     *string `json:"networkName,omitempty"`
	SourceIPAddress *string `json:"sourceIpAddress,omitempty"`
}
