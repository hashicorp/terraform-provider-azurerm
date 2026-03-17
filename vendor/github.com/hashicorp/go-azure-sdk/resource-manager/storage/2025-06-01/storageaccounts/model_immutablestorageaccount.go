package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutableStorageAccount struct {
	Enabled            *bool                                `json:"enabled,omitempty"`
	ImmutabilityPolicy *AccountImmutabilityPolicyProperties `json:"immutabilityPolicy,omitempty"`
}
