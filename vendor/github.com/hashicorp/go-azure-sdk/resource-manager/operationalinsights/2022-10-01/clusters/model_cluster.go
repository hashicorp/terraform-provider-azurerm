package clusters

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Cluster struct {
	Id         *string                                  `json:"id,omitempty"`
	Identity   *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   string                                   `json:"location"`
	Name       *string                                  `json:"name,omitempty"`
	Properties *ClusterProperties                       `json:"properties,omitempty"`
	Sku        *ClusterSku                              `json:"sku,omitempty"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
	Type       *string                                  `json:"type,omitempty"`
}
