package cloudexadatainfrastructures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EstimatedPatchingTime struct {
	EstimatedDbServerPatchingTime        *int64 `json:"estimatedDbServerPatchingTime,omitempty"`
	EstimatedNetworkSwitchesPatchingTime *int64 `json:"estimatedNetworkSwitchesPatchingTime,omitempty"`
	EstimatedStorageServerPatchingTime   *int64 `json:"estimatedStorageServerPatchingTime,omitempty"`
	TotalEstimatedPatchingTime           *int64 `json:"totalEstimatedPatchingTime,omitempty"`
}
