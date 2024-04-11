package configurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationAssignmentProperties struct {
	Filter                     *ConfigurationAssignmentFilterProperties `json:"filter,omitempty"`
	MaintenanceConfigurationId *string                                  `json:"maintenanceConfigurationId,omitempty"`
	ResourceId                 *string                                  `json:"resourceId,omitempty"`
}
