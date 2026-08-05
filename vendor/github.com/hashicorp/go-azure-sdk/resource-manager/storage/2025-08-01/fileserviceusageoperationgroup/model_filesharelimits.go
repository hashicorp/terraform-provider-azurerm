package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareLimits struct {
	MaxProvisionedBandwidthMiBPerSec *int64 `json:"maxProvisionedBandwidthMiBPerSec,omitempty"`
	MaxProvisionedIOPS               *int64 `json:"maxProvisionedIOPS,omitempty"`
	MaxProvisionedStorageGiB         *int64 `json:"maxProvisionedStorageGiB,omitempty"`
	MinProvisionedBandwidthMiBPerSec *int64 `json:"minProvisionedBandwidthMiBPerSec,omitempty"`
	MinProvisionedIOPS               *int64 `json:"minProvisionedIOPS,omitempty"`
	MinProvisionedStorageGiB         *int64 `json:"minProvisionedStorageGiB,omitempty"`
}
