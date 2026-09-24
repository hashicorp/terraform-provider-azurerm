package resourceguardproxy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardOperationDetail struct {
	DefaultResourceRequest *string `json:"defaultResourceRequest,omitempty"`
	VaultCriticalOperation *string `json:"vaultCriticalOperation,omitempty"`
}
