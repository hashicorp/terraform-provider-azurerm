package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RangerUsersyncSpec struct {
	Enabled             *bool               `json:"enabled,omitempty"`
	Groups              *[]string           `json:"groups,omitempty"`
	Mode                *RangerUsersyncMode `json:"mode,omitempty"`
	UserMappingLocation *string             `json:"userMappingLocation,omitempty"`
	Users               *[]string           `json:"users,omitempty"`
}
