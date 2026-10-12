package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProximityPlacementGroup struct {
	EntityTypeIntendedToUse   ProximityPlacementGroupEntityType `json:"entityTypeIntendedToUse"`
	ProximityAnchorId         *string                           `json:"proximityAnchorId,omitempty"`
	ProximityPlacementGroupId string                            `json:"proximityPlacementGroupId"`
}
