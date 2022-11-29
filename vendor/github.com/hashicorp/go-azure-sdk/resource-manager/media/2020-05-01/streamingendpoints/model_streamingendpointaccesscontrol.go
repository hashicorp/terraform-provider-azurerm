package streamingendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingEndpointAccessControl struct {
	Akamai *AkamaiAccessControl `json:"akamai"`
	IP     *IPAccessControl     `json:"ip"`
}
