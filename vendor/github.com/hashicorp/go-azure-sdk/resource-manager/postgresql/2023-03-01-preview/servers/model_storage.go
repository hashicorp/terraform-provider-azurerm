package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Storage struct {
	AutoGrow      *StorageAutoGrow                  `json:"autoGrow,omitempty"`
	Iops          *int64                            `json:"iops,omitempty"`
	StorageSizeGB *int64                            `json:"storageSizeGB,omitempty"`
	Tier          *AzureManagedDiskPerformanceTiers `json:"tier,omitempty"`
}
