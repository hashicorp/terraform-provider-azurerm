package saprecommendations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSizingRecommendationRequest struct {
	AppLocation          string                   `json:"appLocation"`
	DatabaseType         SAPDatabaseType          `json:"databaseType"`
	DbMemory             int64                    `json:"dbMemory"`
	DbScaleMethod        *SAPDatabaseScaleMethod  `json:"dbScaleMethod,omitempty"`
	DeploymentType       SAPDeploymentType        `json:"deploymentType"`
	Environment          SAPEnvironmentType       `json:"environment"`
	HighAvailabilityType *SAPHighAvailabilityType `json:"highAvailabilityType,omitempty"`
	SapProduct           SAPProductType           `json:"sapProduct"`
	Saps                 int64                    `json:"saps"`
}
