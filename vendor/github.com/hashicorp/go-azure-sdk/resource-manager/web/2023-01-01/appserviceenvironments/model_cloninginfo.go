package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloningInfo struct {
	AppSettingsOverrides      *map[string]string `json:"appSettingsOverrides,omitempty"`
	CloneCustomHostNames      *bool              `json:"cloneCustomHostNames,omitempty"`
	CloneSourceControl        *bool              `json:"cloneSourceControl,omitempty"`
	ConfigureLoadBalancing    *bool              `json:"configureLoadBalancing,omitempty"`
	CorrelationId             *string            `json:"correlationId,omitempty"`
	HostingEnvironment        *string            `json:"hostingEnvironment,omitempty"`
	Overwrite                 *bool              `json:"overwrite,omitempty"`
	SourceWebAppId            string             `json:"sourceWebAppId"`
	SourceWebAppLocation      *string            `json:"sourceWebAppLocation,omitempty"`
	TrafficManagerProfileId   *string            `json:"trafficManagerProfileId,omitempty"`
	TrafficManagerProfileName *string            `json:"trafficManagerProfileName,omitempty"`
}
