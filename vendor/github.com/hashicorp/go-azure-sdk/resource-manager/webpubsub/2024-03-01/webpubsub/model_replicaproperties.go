package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicaProperties struct {
	ProvisioningState     *ProvisioningState `json:"provisioningState,omitempty"`
	RegionEndpointEnabled *string            `json:"regionEndpointEnabled,omitempty"`
	ResourceStopped       *string            `json:"resourceStopped,omitempty"`
}
