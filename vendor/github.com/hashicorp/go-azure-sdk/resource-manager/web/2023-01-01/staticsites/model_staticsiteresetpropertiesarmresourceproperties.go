package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteResetPropertiesARMResourceProperties struct {
	RepositoryToken        *string `json:"repositoryToken,omitempty"`
	ShouldUpdateRepository *bool   `json:"shouldUpdateRepository,omitempty"`
}
