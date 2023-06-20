package videoanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzerIdentity struct {
	Type                   string                                  `json:"type"`
	UserAssignedIdentities *map[string]UserAssignedManagedIdentity `json:"userAssignedIdentities,omitempty"`
}
