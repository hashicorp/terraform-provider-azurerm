package projects

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InheritedSettingsForProject struct {
	NetworkSettings        *ProjectNetworkSettings          `json:"networkSettings,omitempty"`
	ProjectCatalogSettings *DevCenterProjectCatalogSettings `json:"projectCatalogSettings,omitempty"`
}
