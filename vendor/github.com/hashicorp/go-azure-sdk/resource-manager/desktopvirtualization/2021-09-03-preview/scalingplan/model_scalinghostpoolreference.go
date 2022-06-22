package scalingplan

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScalingHostPoolReference struct {
	HostPoolArmPath    *string `json:"hostPoolArmPath,omitempty"`
	ScalingPlanEnabled *bool   `json:"scalingPlanEnabled,omitempty"`
}
