package servervalidateestimatehighavailability

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HighAvailabilityValidationEstimation struct {
	EstimatedDowntime                *int64  `json:"estimatedDowntime,omitempty"`
	ExpectedStandbyAvailabilityZone  *string `json:"expectedStandbyAvailabilityZone,omitempty"`
	ScheduledStandbyAvailabilityZone *string `json:"scheduledStandbyAvailabilityZone,omitempty"`
}
