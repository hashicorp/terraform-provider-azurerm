package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CircuitBreakerRule struct {
	AcceptRetryAfter *bool                           `json:"acceptRetryAfter,omitempty"`
	FailureCondition *CircuitBreakerFailureCondition `json:"failureCondition,omitempty"`
	Name             *string                         `json:"name,omitempty"`
	TripDuration     *string                         `json:"tripDuration,omitempty"`
}
