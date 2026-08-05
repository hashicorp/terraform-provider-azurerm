package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeRelocationProperties struct {
	ReadyToBeFinalized  *bool `json:"readyToBeFinalized,omitempty"`
	RelocationRequested *bool `json:"relocationRequested,omitempty"`
}
