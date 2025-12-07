package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceTypeSchemaResources struct {
	Limits   *map[string]string `json:"limits,omitempty"`
	Requests *map[string]string `json:"requests,omitempty"`
}
