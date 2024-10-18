package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MobilityServiceUpdate struct {
	OsType       *string `json:"osType,omitempty"`
	RebootStatus *string `json:"rebootStatus,omitempty"`
	Version      *string `json:"version,omitempty"`
}
