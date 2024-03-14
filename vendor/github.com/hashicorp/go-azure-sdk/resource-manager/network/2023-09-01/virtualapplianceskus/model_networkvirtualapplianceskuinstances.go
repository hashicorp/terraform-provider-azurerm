package virtualapplianceskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualApplianceSkuInstances struct {
	InstanceCount *int64  `json:"instanceCount,omitempty"`
	ScaleUnit     *string `json:"scaleUnit,omitempty"`
}
