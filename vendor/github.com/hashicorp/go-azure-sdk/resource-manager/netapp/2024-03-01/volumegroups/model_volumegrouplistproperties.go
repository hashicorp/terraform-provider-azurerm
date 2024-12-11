package volumegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupListProperties struct {
	GroupMetaData     *VolumeGroupMetaData `json:"groupMetaData,omitempty"`
	ProvisioningState *string              `json:"provisioningState,omitempty"`
}
