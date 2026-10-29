package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolPlatformProfile struct {
	AvailabilityZone       *string        `json:"availabilityZone,omitempty"`
	EnableEncryptionAtHost *bool          `json:"enableEncryptionAtHost,omitempty"`
	OsDisk                 *OsDiskProfile `json:"osDisk,omitempty"`
	SubnetId               *string        `json:"subnetId,omitempty"`
	VMSize                 string         `json:"vmSize"`
}
