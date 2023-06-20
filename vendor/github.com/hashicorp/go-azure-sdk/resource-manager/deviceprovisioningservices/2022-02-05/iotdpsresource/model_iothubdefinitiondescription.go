package iotdpsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotHubDefinitionDescription struct {
	AllocationWeight      *int64  `json:"allocationWeight,omitempty"`
	ApplyAllocationPolicy *bool   `json:"applyAllocationPolicy,omitempty"`
	ConnectionString      string  `json:"connectionString"`
	Location              string  `json:"location"`
	Name                  *string `json:"name,omitempty"`
}
