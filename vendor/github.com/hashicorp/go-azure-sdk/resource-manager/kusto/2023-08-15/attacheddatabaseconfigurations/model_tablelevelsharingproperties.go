package attacheddatabaseconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableLevelSharingProperties struct {
	ExternalTablesToExclude    *[]string `json:"externalTablesToExclude,omitempty"`
	ExternalTablesToInclude    *[]string `json:"externalTablesToInclude,omitempty"`
	FunctionsToExclude         *[]string `json:"functionsToExclude,omitempty"`
	FunctionsToInclude         *[]string `json:"functionsToInclude,omitempty"`
	MaterializedViewsToExclude *[]string `json:"materializedViewsToExclude,omitempty"`
	MaterializedViewsToInclude *[]string `json:"materializedViewsToInclude,omitempty"`
	TablesToExclude            *[]string `json:"tablesToExclude,omitempty"`
	TablesToInclude            *[]string `json:"tablesToInclude,omitempty"`
}
