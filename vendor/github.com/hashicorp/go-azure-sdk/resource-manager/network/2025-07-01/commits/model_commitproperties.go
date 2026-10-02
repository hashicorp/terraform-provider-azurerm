package commits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommitProperties struct {
	ActiveLocations   *[]string          `json:"activeLocations,omitempty"`
	CommitType        ConfigurationType  `json:"commitType"`
	ConfigurationIds  *[]string          `json:"configurationIds,omitempty"`
	Description       *string            `json:"description,omitempty"`
	ForceUpdateTag    *string            `json:"forceUpdateTag,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceGuid      *string            `json:"resourceGuid,omitempty"`
	TargetLocations   []string           `json:"targetLocations"`
}
