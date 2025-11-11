package softwareupdateconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetProperties struct {
	AzureQueries    *[]AzureQueryProperties    `json:"azureQueries,omitempty"`
	NonAzureQueries *[]NonAzureQueryProperties `json:"nonAzureQueries,omitempty"`
}
