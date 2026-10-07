package catalogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncErrorDetails struct {
	Conflicts      *[]CatalogConflictError `json:"conflicts,omitempty"`
	Errors         *[]CatalogSyncError     `json:"errors,omitempty"`
	OperationError *CatalogErrorDetails    `json:"operationError,omitempty"`
}
