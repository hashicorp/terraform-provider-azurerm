package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecuritySettings struct {
	ImmutabilitySettings   *ImmutabilitySettings   `json:"immutabilitySettings,omitempty"`
	MultiUserAuthorization *MultiUserAuthorization `json:"multiUserAuthorization,omitempty"`
	SoftDeleteSettings     *SoftDeleteSettings     `json:"softDeleteSettings,omitempty"`
}
