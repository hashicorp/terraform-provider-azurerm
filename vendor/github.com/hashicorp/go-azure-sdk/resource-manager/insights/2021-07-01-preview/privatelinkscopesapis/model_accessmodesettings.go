package privatelinkscopesapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessModeSettings struct {
	Exclusions          *[]AccessModeSettingsExclusion `json:"exclusions,omitempty"`
	IngestionAccessMode AccessMode                     `json:"ingestionAccessMode"`
	QueryAccessMode     AccessMode                     `json:"queryAccessMode"`
}
