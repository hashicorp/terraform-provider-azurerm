package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerNetworkInterfaceStatistics struct {
	RxBytes   *int64 `json:"rxBytes,omitempty"`
	RxDropped *int64 `json:"rxDropped,omitempty"`
	RxErrors  *int64 `json:"rxErrors,omitempty"`
	RxPackets *int64 `json:"rxPackets,omitempty"`
	TxBytes   *int64 `json:"txBytes,omitempty"`
	TxDropped *int64 `json:"txDropped,omitempty"`
	TxErrors  *int64 `json:"txErrors,omitempty"`
	TxPackets *int64 `json:"txPackets,omitempty"`
}
