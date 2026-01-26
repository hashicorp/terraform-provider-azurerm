package securitydomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityDomainOperationStatus struct {
	Status        *OperationStatus `json:"status,omitempty"`
	StatusDetails *string          `json:"status_details,omitempty"`
}
