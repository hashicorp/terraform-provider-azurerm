package backupinstanceresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceListSelectionCriteria struct {
	ObjectType            string             `json:"objectType"`
	ResourceIdentifiers   []string           `json:"resourceIdentifiers"`
	ResourceNameOverrides *map[string]string `json:"resourceNameOverrides,omitempty"`
}
