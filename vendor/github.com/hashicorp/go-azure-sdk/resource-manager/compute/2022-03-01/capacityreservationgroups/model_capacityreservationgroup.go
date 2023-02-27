package capacityreservationgroups

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationGroup struct {
	Id         *string                             `json:"id,omitempty"`
	Location   string                              `json:"location"`
	Name       *string                             `json:"name,omitempty"`
	Properties *CapacityReservationGroupProperties `json:"properties,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
	Type       *string                             `json:"type,omitempty"`
	Zones      *zones.Schema                       `json:"zones,omitempty"`
}
