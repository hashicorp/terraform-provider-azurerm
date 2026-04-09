package blobcontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutabilityPolicyProperty struct {
	AllowProtectedAppendWrites            *bool                    `json:"allowProtectedAppendWrites,omitempty"`
	AllowProtectedAppendWritesAll         *bool                    `json:"allowProtectedAppendWritesAll,omitempty"`
	ImmutabilityPeriodSinceCreationInDays *int64                   `json:"immutabilityPeriodSinceCreationInDays,omitempty"`
	State                                 *ImmutabilityPolicyState `json:"state,omitempty"`
}
