package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceTypeSchema struct {
	NodeSelector *map[string]string           `json:"nodeSelector,omitempty"`
	Resources    *InstanceTypeSchemaResources `json:"resources,omitempty"`
}
