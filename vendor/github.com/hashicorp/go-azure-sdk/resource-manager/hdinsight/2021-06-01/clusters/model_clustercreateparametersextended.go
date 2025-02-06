package clusters

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterCreateParametersExtended struct {
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Properties *ClusterCreateProperties           `json:"properties,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Zones      *zones.Schema                      `json:"zones,omitempty"`
}
