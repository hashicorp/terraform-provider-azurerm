package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlinkStorageProfile struct {
	StorageUri string  `json:"storageUri"`
	Storagekey *string `json:"storagekey,omitempty"`
}
