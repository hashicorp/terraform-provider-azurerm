package blobcontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LeaseContainerRequest struct {
	Action          LeaseContainerRequestAction `json:"action"`
	BreakPeriod     *int64                      `json:"breakPeriod,omitempty"`
	LeaseDuration   *int64                      `json:"leaseDuration,omitempty"`
	LeaseId         *string                     `json:"leaseId,omitempty"`
	ProposedLeaseId *string                     `json:"proposedLeaseId,omitempty"`
}
