package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountImmutabilityPolicyProperties struct {
	AllowProtectedAppendWrites            *bool                           `json:"allowProtectedAppendWrites,omitempty"`
	ImmutabilityPeriodSinceCreationInDays *int64                          `json:"immutabilityPeriodSinceCreationInDays,omitempty"`
	State                                 *AccountImmutabilityPolicyState `json:"state,omitempty"`
}
