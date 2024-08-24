package managedclusters

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCluster struct {
	ExtendedLocation *edgezones.Model                  `json:"extendedLocation,omitempty"`
	Id               *string                           `json:"id,omitempty"`
	Identity         *identity.SystemOrUserAssignedMap `json:"identity,omitempty"`
	Location         string                            `json:"location"`
	Name             *string                           `json:"name,omitempty"`
	Properties       *ManagedClusterProperties         `json:"properties,omitempty"`
	Sku              *ManagedClusterSKU                `json:"sku,omitempty"`
	SystemData       *systemdata.SystemData            `json:"systemData,omitempty"`
	Tags             *map[string]string                `json:"tags,omitempty"`
	Type             *string                           `json:"type,omitempty"`
}
