package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkerProfile struct {
	Count               *int64            `json:"count,omitempty"`
	DiskEncryptionSetId *string           `json:"diskEncryptionSetId,omitempty"`
	DiskSizeGB          *int64            `json:"diskSizeGB,omitempty"`
	EncryptionAtHost    *EncryptionAtHost `json:"encryptionAtHost,omitempty"`
	Name                *string           `json:"name,omitempty"`
	SubnetId            *string           `json:"subnetId,omitempty"`
	VMSize              *string           `json:"vmSize,omitempty"`
}
