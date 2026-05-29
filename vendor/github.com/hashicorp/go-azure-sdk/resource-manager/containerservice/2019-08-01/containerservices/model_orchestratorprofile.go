package containerservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrchestratorProfile struct {
	IsPreview           *bool   `json:"isPreview,omitempty"`
	OrchestratorType    *string `json:"orchestratorType,omitempty"`
	OrchestratorVersion string  `json:"orchestratorVersion"`
}
