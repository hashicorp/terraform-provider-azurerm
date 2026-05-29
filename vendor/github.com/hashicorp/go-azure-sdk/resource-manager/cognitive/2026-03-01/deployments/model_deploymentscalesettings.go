package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScaleSettings struct {
	ActiveCapacity *int64               `json:"activeCapacity,omitempty"`
	Capacity       *int64               `json:"capacity,omitempty"`
	ScaleType      *DeploymentScaleType `json:"scaleType,omitempty"`
}
