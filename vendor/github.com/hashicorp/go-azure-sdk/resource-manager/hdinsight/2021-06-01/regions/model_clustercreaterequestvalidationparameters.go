package regions

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterCreateRequestValidationParameters struct {
	FetchAaddsResource *bool                              `json:"fetchAaddsResource,omitempty"`
	Identity           *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location           *string                            `json:"location,omitempty"`
	Name               *string                            `json:"name,omitempty"`
	Properties         *ClusterCreateProperties           `json:"properties,omitempty"`
	Tags               *map[string]string                 `json:"tags,omitempty"`
	TenantId           *string                            `json:"tenantId,omitempty"`
	Type               *string                            `json:"type,omitempty"`
	Zones              *zones.Schema                      `json:"zones,omitempty"`
}
