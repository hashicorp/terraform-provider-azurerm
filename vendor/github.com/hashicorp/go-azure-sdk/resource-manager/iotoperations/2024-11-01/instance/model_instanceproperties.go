package instance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceProperties struct {
	Description       *string            `json:"description,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	SchemaRegistryRef SchemaRegistryRef  `json:"schemaRegistryRef"`
	Version           *string            `json:"version,omitempty"`
}
