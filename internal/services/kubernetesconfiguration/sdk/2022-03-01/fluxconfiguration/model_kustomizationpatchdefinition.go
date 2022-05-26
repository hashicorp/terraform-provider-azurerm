package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KustomizationPatchDefinition struct {
	DependsOn              *[]string `json:"dependsOn,omitempty"`
	Force                  *bool     `json:"force,omitempty"`
	Path                   *string   `json:"path,omitempty"`
	Prune                  *bool     `json:"prune,omitempty"`
	RetryIntervalInSeconds *int64    `json:"retryIntervalInSeconds,omitempty"`
	SyncIntervalInSeconds  *int64    `json:"syncIntervalInSeconds,omitempty"`
	TimeoutInSeconds       *int64    `json:"timeoutInSeconds,omitempty"`
}
