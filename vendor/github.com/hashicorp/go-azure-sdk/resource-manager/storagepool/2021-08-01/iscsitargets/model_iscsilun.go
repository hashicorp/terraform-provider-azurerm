package iscsitargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IscsiLun struct {
	Lun                        *int64 `json:"lun,omitempty"`
	ManagedDiskAzureResourceId string `json:"managedDiskAzureResourceId"`
	Name                       string `json:"name"`
}
