package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAccessRequestModel struct {
	SearchFilters *map[string]string `json:"searchFilters,omitempty"`
}
