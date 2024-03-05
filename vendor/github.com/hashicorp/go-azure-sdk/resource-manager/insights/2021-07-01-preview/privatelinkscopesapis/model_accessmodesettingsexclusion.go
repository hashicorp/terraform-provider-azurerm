package privatelinkscopesapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessModeSettingsExclusion struct {
	IngestionAccessMode           *AccessMode `json:"ingestionAccessMode,omitempty"`
	PrivateEndpointConnectionName *string     `json:"privateEndpointConnectionName,omitempty"`
	QueryAccessMode               *AccessMode `json:"queryAccessMode,omitempty"`
}
