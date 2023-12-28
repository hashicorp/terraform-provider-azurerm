package packetcaptures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureResultProperties struct {
	BytesToCapturePerPacket *int64                       `json:"bytesToCapturePerPacket,omitempty"`
	Filters                 *[]PacketCaptureFilter       `json:"filters,omitempty"`
	ProvisioningState       *ProvisioningState           `json:"provisioningState,omitempty"`
	Scope                   *PacketCaptureMachineScope   `json:"scope,omitempty"`
	StorageLocation         PacketCaptureStorageLocation `json:"storageLocation"`
	Target                  string                       `json:"target"`
	TargetType              *PacketCaptureTargetType     `json:"targetType,omitempty"`
	TimeLimitInSeconds      *int64                       `json:"timeLimitInSeconds,omitempty"`
	TotalBytesPerSession    *int64                       `json:"totalBytesPerSession,omitempty"`
}
