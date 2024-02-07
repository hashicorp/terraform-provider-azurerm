package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceConfig struct {
	Cpu    *float64 `json:"cpu,omitempty"`
	Memory *string  `json:"memory,omitempty"`
}
