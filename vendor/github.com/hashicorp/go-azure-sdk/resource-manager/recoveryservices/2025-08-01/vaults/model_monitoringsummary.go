package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoringSummary struct {
	DeprecatedProviderCount  *int64 `json:"deprecatedProviderCount,omitempty"`
	EventsCount              *int64 `json:"eventsCount,omitempty"`
	SupportedProviderCount   *int64 `json:"supportedProviderCount,omitempty"`
	UnHealthyProviderCount   *int64 `json:"unHealthyProviderCount,omitempty"`
	UnHealthyVMCount         *int64 `json:"unHealthyVmCount,omitempty"`
	UnsupportedProviderCount *int64 `json:"unsupportedProviderCount,omitempty"`
}
