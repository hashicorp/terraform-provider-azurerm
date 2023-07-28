package batchdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceConfiguration struct {
	InstanceCount *int64                  `json:"instanceCount,omitempty"`
	InstanceType  *string                 `json:"instanceType,omitempty"`
	Properties    *map[string]interface{} `json:"properties,omitempty"`
}
