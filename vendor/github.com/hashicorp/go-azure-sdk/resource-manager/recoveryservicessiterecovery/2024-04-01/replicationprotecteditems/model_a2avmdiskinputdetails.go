package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AVMDiskInputDetails struct {
	DiskUri                             string `json:"diskUri"`
	PrimaryStagingAzureStorageAccountId string `json:"primaryStagingAzureStorageAccountId"`
	RecoveryAzureStorageAccountId       string `json:"recoveryAzureStorageAccountId"`
}
