package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlatformUpdateParameters struct {
	Architecture *Architecture `json:"architecture,omitempty"`
	Os           *OS           `json:"os,omitempty"`
	Variant      *Variant      `json:"variant,omitempty"`
}
