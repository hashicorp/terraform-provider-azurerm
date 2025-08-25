package containerapps

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerApp struct {
	ExtendedLocation *ExtendedLocation                        `json:"extendedLocation,omitempty"`
	Id               *string                                  `json:"id,omitempty"`
	Identity         *identity.LegacySystemAndUserAssignedMap `json:"identity,omitempty"`
	Location         string                                   `json:"location"`
	ManagedBy        *string                                  `json:"managedBy,omitempty"`
	Name             *string                                  `json:"name,omitempty"`
	Properties       *ContainerAppProperties                  `json:"properties,omitempty"`
	SystemData       *systemdata.SystemData                   `json:"systemData,omitempty"`
	Tags             *map[string]string                       `json:"tags,omitempty"`
	Type             *string                                  `json:"type,omitempty"`
}
