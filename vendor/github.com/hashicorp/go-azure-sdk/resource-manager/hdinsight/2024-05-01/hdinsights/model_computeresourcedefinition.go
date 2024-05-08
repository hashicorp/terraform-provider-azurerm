package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeResourceDefinition struct {
	Cpu    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
}
