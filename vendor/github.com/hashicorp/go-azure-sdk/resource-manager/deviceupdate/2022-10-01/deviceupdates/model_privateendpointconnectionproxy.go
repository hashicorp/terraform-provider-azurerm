package deviceupdates

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointConnectionProxy struct {
	ETag                  *string                                   `json:"eTag,omitempty"`
	Id                    *string                                   `json:"id,omitempty"`
	Name                  *string                                   `json:"name,omitempty"`
	Properties            *PrivateEndpointConnectionProxyProperties `json:"properties,omitempty"`
	RemotePrivateEndpoint *RemotePrivateEndpoint                    `json:"remotePrivateEndpoint,omitempty"`
	Status                *string                                   `json:"status,omitempty"`
	SystemData            *systemdata.SystemData                    `json:"systemData,omitempty"`
	Type                  *string                                   `json:"type,omitempty"`
}
