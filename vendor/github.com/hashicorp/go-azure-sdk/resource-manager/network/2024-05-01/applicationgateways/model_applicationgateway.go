package applicationgateways

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGateway struct {
	Etag       *string                             `json:"etag,omitempty"`
	Id         *string                             `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap  `json:"identity,omitempty"`
	Location   *string                             `json:"location,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *ApplicationGatewayPropertiesFormat `json:"properties,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
	Type       *string                             `json:"type,omitempty"`
	Zones      *zones.Schema                       `json:"zones,omitempty"`
}
