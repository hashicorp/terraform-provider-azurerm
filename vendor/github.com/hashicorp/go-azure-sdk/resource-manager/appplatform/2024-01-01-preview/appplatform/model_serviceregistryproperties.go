package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceRegistryProperties struct {
	Instances         *[]ServiceRegistryInstance        `json:"instances,omitempty"`
	ProvisioningState *ServiceRegistryProvisioningState `json:"provisioningState,omitempty"`
	ResourceRequests  *ServiceRegistryResourceRequests  `json:"resourceRequests,omitempty"`
}
