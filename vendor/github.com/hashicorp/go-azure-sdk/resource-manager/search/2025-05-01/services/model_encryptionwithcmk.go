package services

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionWithCmk struct {
	EncryptionComplianceStatus *SearchEncryptionComplianceStatus `json:"encryptionComplianceStatus,omitempty"`
	Enforcement                *SearchEncryptionWithCmk          `json:"enforcement,omitempty"`
}
