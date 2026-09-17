package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsScaleAndConcurrency struct {
	AlwaysReady          *[]FunctionsAlwaysReadyConfig         `json:"alwaysReady,omitempty"`
	InstanceMemoryMB     *int64                                `json:"instanceMemoryMB,omitempty"`
	MaximumInstanceCount *int64                                `json:"maximumInstanceCount,omitempty"`
	Triggers             *FunctionsScaleAndConcurrencyTriggers `json:"triggers,omitempty"`
}
