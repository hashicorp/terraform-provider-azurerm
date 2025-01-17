package virtualnetworkgateways

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGateway struct {
	Etag             *string                               `json:"etag,omitempty"`
	ExtendedLocation *edgezones.Model                      `json:"extendedLocation,omitempty"`
	Id               *string                               `json:"id,omitempty"`
	Identity         *identity.SystemAndUserAssignedMap    `json:"identity,omitempty"`
	Location         *string                               `json:"location,omitempty"`
	Name             *string                               `json:"name,omitempty"`
	Properties       VirtualNetworkGatewayPropertiesFormat `json:"properties"`
	Tags             *map[string]string                    `json:"tags,omitempty"`
	Type             *string                               `json:"type,omitempty"`
}
