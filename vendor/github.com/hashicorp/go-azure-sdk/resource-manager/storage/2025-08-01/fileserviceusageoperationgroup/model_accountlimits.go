package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountLimits struct {
	MaxFileShares                    *int64 `json:"maxFileShares,omitempty"`
	MaxProvisionedBandwidthMiBPerSec *int64 `json:"maxProvisionedBandwidthMiBPerSec,omitempty"`
	MaxProvisionedIOPS               *int64 `json:"maxProvisionedIOPS,omitempty"`
	MaxProvisionedStorageGiB         *int64 `json:"maxProvisionedStorageGiB,omitempty"`
}
