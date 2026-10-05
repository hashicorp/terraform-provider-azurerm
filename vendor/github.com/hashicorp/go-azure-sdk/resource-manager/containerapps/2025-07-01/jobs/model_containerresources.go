package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerResources struct {
	Cpu              *float64 `json:"cpu,omitempty"`
	EphemeralStorage *string  `json:"ephemeralStorage,omitempty"`
	Memory           *string  `json:"memory,omitempty"`
}
