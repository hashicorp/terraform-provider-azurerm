package eventhubsclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	CreatedAt         *string            `json:"createdAt,omitempty"`
	MetricId          *string            `json:"metricId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Status            *string            `json:"status,omitempty"`
	SupportsScaling   *bool              `json:"supportsScaling,omitempty"`
	UpdatedAt         *string            `json:"updatedAt,omitempty"`
}
