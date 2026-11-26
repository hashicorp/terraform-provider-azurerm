package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolPropertiesForPatchUpdate struct {
	ComputeModel         *PoolComputeModelForPatchUpdate               `json:"computeModel,omitempty"`
	HighAvailability     *PoolHighAvailabilityPatch                    `json:"highAvailability,omitempty"`
	NodeCount            *int64                                        `json:"nodeCount,omitempty"`
	ParameterGroup       *PoolParameterGroupReferenceForPatchUpdate    `json:"parameterGroup,omitempty"`
	ScheduledMaintenance *ScheduledMaintenancePropertiesForPatchUpdate `json:"scheduledMaintenance,omitempty"`
}
