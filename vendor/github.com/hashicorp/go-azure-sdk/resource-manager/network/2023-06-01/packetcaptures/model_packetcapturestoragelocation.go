package packetcaptures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureStorageLocation struct {
	FilePath    *string `json:"filePath,omitempty"`
	StorageId   *string `json:"storageId,omitempty"`
	StoragePath *string `json:"storagePath,omitempty"`
}
