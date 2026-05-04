package clouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageQosPolicy struct {
	BandwidthLimit *int64  `json:"bandwidthLimit,omitempty"`
	Id             *string `json:"id,omitempty"`
	IopsMaximum    *int64  `json:"iopsMaximum,omitempty"`
	IopsMinimum    *int64  `json:"iopsMinimum,omitempty"`
	Name           *string `json:"name,omitempty"`
	PolicyId       *string `json:"policyId,omitempty"`
}
