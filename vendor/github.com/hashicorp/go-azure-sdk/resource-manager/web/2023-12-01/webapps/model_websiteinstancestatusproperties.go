package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebSiteInstanceStatusProperties struct {
	ConsoleUrl     *string                   `json:"consoleUrl,omitempty"`
	Containers     *map[string]ContainerInfo `json:"containers,omitempty"`
	DetectorUrl    *string                   `json:"detectorUrl,omitempty"`
	HealthCheckUrl *string                   `json:"healthCheckUrl,omitempty"`
	State          *SiteRuntimeState         `json:"state,omitempty"`
	StatusUrl      *string                   `json:"statusUrl,omitempty"`
}
