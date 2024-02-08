package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentLocations struct {
	HostingEnvironmentDeploymentInfos *[]HostingEnvironmentDeploymentInfo `json:"hostingEnvironmentDeploymentInfos,omitempty"`
	HostingEnvironments               *[]AppServiceEnvironment            `json:"hostingEnvironments,omitempty"`
	Locations                         *[]GeoRegion                        `json:"locations,omitempty"`
}
