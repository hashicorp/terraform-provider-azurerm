package catalogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncStats struct {
	Added                  *int64             `json:"added,omitempty"`
	Removed                *int64             `json:"removed,omitempty"`
	SyncedCatalogItemTypes *[]CatalogItemType `json:"syncedCatalogItemTypes,omitempty"`
	SynchronizationErrors  *int64             `json:"synchronizationErrors,omitempty"`
	Unchanged              *int64             `json:"unchanged,omitempty"`
	Updated                *int64             `json:"updated,omitempty"`
	ValidationErrors       *int64             `json:"validationErrors,omitempty"`
}
