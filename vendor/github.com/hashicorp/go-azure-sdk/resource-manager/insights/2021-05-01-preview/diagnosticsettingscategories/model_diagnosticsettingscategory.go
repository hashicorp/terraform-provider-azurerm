package diagnosticsettingscategories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsCategory struct {
	CategoryGroups *[]string     `json:"categoryGroups,omitempty"`
	CategoryType   *CategoryType `json:"categoryType,omitempty"`
}
