package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceForestSettings struct {
	ResourceForest *string        `json:"resourceForest,omitempty"`
	Settings       *[]ForestTrust `json:"settings,omitempty"`
}
