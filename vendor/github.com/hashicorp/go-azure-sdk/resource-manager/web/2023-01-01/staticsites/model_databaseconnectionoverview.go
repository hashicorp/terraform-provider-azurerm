package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseConnectionOverview struct {
	ConfigurationFiles *[]StaticSiteDatabaseConnectionConfigurationFileOverview `json:"configurationFiles,omitempty"`
	ConnectionIdentity *string                                                  `json:"connectionIdentity,omitempty"`
	Name               *string                                                  `json:"name,omitempty"`
	Region             *string                                                  `json:"region,omitempty"`
	ResourceId         *string                                                  `json:"resourceId,omitempty"`
}
