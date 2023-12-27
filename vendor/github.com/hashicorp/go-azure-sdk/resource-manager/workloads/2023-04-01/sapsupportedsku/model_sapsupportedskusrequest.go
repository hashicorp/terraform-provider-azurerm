package sapsupportedsku

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSupportedSkusRequest struct {
	AppLocation          string                   `json:"appLocation"`
	DatabaseType         SAPDatabaseType          `json:"databaseType"`
	DeploymentType       SAPDeploymentType        `json:"deploymentType"`
	Environment          SAPEnvironmentType       `json:"environment"`
	HighAvailabilityType *SAPHighAvailabilityType `json:"highAvailabilityType,omitempty"`
	SapProduct           SAPProductType           `json:"sapProduct"`
}
