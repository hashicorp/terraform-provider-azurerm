package resourceguards

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardOperation struct {
	RequestResourceType    *string `json:"requestResourceType,omitempty"`
	VaultCriticalOperation *string `json:"vaultCriticalOperation,omitempty"`
}
