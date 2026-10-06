package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebSiteInstanceStatusProperties struct {
	ConsoleURL     *string                   `json:"consoleUrl,omitempty"`
	Containers     *map[string]ContainerInfo `json:"containers,omitempty"`
	DetectorURL    *string                   `json:"detectorUrl,omitempty"`
	HealthCheckURL *string                   `json:"healthCheckUrl,omitempty"`
	State          *SiteRuntimeState         `json:"state,omitempty"`
	StatusURL      *string                   `json:"statusUrl,omitempty"`
}
