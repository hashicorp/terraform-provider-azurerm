package resourcemanagementprivatelink

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceManagementPrivateLink struct {
	Id         *string                                           `json:"id,omitempty"`
	Location   *string                                           `json:"location,omitempty"`
	Name       *string                                           `json:"name,omitempty"`
	Properties *ResourceManagementPrivateLinkEndpointConnections `json:"properties,omitempty"`
	Type       *string                                           `json:"type,omitempty"`
}
