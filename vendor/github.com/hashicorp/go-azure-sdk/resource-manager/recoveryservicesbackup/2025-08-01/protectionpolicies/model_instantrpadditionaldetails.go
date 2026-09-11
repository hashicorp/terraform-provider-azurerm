package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstantRPAdditionalDetails struct {
	AzureBackupRGNamePrefix *string `json:"azureBackupRGNamePrefix,omitempty"`
	AzureBackupRGNameSuffix *string `json:"azureBackupRGNameSuffix,omitempty"`
}
