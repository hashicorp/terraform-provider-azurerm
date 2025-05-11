package catalogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CatalogUpdateProperties struct {
	AdoGit   *GitCatalog        `json:"adoGit,omitempty"`
	GitHub   *GitCatalog        `json:"gitHub,omitempty"`
	SyncType *CatalogSyncType   `json:"syncType,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
}
