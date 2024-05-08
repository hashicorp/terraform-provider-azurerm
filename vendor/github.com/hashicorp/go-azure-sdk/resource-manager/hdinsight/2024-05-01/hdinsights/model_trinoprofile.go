package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrinoProfile struct {
	CatalogOptions    *CatalogOptions     `json:"catalogOptions,omitempty"`
	Coordinator       *TrinoCoordinator   `json:"coordinator,omitempty"`
	UserPluginsSpec   *TrinoUserPlugins   `json:"userPluginsSpec,omitempty"`
	UserTelemetrySpec *TrinoUserTelemetry `json:"userTelemetrySpec,omitempty"`
	Worker            *TrinoWorker        `json:"worker,omitempty"`
}
