package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MeshRevision struct {
	CompatibleWith *[]CompatibleVersions `json:"compatibleWith,omitempty"`
	Revision       *string               `json:"revision,omitempty"`
	Upgrades       *[]string             `json:"upgrades,omitempty"`
}
