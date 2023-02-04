package containerservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrchestratorVersionProfile struct {
	Default             *bool                  `json:"default,omitempty"`
	IsPreview           *bool                  `json:"isPreview,omitempty"`
	OrchestratorType    string                 `json:"orchestratorType"`
	OrchestratorVersion string                 `json:"orchestratorVersion"`
	Upgrades            *[]OrchestratorProfile `json:"upgrades,omitempty"`
}
