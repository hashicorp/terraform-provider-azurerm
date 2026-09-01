package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeMonitoringData struct {
	Name  *string                                 `json:"name,omitempty"`
	Nodes *[]IntegrationRuntimeNodeMonitoringData `json:"nodes,omitempty"`
}
