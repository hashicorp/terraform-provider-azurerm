package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutionPlacement struct {
	Constraints  *[]PlacementConstraint `json:"constraints,omitempty"`
	Distribution *DistributionPolicy    `json:"distribution,omitempty"`
}
