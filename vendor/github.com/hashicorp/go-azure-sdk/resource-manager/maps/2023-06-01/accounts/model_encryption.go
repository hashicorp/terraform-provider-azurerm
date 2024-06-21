package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Encryption struct {
	CustomerManagedKeyEncryption *CustomerManagedKeyEncryption `json:"customerManagedKeyEncryption,omitempty"`
	InfrastructureEncryption     *InfrastructureEncryption     `json:"infrastructureEncryption,omitempty"`
}
