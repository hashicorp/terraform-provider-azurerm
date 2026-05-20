package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Storage struct {
	AutoGrow      *EnableStatusEnum `json:"autoGrow,omitempty"`
	AutoIoScaling *EnableStatusEnum `json:"autoIoScaling,omitempty"`
	Iops          *int64            `json:"iops,omitempty"`
	LogOnDisk     *EnableStatusEnum `json:"logOnDisk,omitempty"`
	StorageSizeGB *int64            `json:"storageSizeGB,omitempty"`
	StorageSku    *string           `json:"storageSku,omitempty"`
}
