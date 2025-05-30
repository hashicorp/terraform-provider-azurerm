package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisMigrationInfo struct {
	EnvironmentOverwriteOption *SsisMigrationOverwriteOption `json:"environmentOverwriteOption,omitempty"`
	ProjectOverwriteOption     *SsisMigrationOverwriteOption `json:"projectOverwriteOption,omitempty"`
	SsisStoreType              *SsisStoreType                `json:"ssisStoreType,omitempty"`
}
