package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskBackedMessageBuffer struct {
	EphemeralVolumeClaimSpec  *VolumeClaimSpec `json:"ephemeralVolumeClaimSpec,omitempty"`
	MaxSize                   string           `json:"maxSize"`
	PersistentVolumeClaimSpec *VolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
}
