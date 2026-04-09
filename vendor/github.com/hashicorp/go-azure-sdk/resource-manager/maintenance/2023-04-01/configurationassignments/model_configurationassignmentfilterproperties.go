package configurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationAssignmentFilterProperties struct {
	Locations      *[]string              `json:"locations,omitempty"`
	OsTypes        *[]string              `json:"osTypes,omitempty"`
	ResourceGroups *[]string              `json:"resourceGroups,omitempty"`
	ResourceTypes  *[]string              `json:"resourceTypes,omitempty"`
	TagSettings    *TagSettingsProperties `json:"tagSettings,omitempty"`
}
