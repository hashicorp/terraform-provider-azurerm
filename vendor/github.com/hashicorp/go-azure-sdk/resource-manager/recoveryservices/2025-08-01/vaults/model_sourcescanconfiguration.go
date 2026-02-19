package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceScanConfiguration struct {
	SourceScanIdentity *AssociatedIdentity `json:"sourceScanIdentity,omitempty"`
	State              *State              `json:"state,omitempty"`
}
