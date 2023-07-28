package environmentcontainer

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentContainer struct {
	Description       *string                 `json:"description,omitempty"`
	IsArchived        *bool                   `json:"isArchived,omitempty"`
	LatestVersion     *string                 `json:"latestVersion,omitempty"`
	NextVersion       *string                 `json:"nextVersion,omitempty"`
	Properties        *map[string]string      `json:"properties,omitempty"`
	ProvisioningState *AssetProvisioningState `json:"provisioningState,omitempty"`
	Tags              *map[string]string      `json:"tags,omitempty"`
}
