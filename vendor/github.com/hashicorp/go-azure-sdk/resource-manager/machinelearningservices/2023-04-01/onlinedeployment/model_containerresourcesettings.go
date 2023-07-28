package onlinedeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerResourceSettings struct {
	Cpu    *string `json:"cpu,omitempty"`
	Gpu    *string `json:"gpu,omitempty"`
	Memory *string `json:"memory,omitempty"`
}
