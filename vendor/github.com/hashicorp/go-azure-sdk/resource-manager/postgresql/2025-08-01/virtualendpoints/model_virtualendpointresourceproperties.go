package virtualendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualEndpointResourceProperties struct {
	EndpointType     *VirtualEndpointType `json:"endpointType,omitempty"`
	Members          *[]string            `json:"members,omitempty"`
	VirtualEndpoints *[]string            `json:"virtualEndpoints,omitempty"`
}
