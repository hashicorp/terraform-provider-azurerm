package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecuteSSISPackageActivityTypeProperties struct {
	ConnectVia                IntegrationRuntimeReference                   `json:"connectVia"`
	EnvironmentPath           *interface{}                                  `json:"environmentPath,omitempty"`
	ExecutionCredential       *SSISExecutionCredential                      `json:"executionCredential,omitempty"`
	LogLocation               *SSISLogLocation                              `json:"logLocation,omitempty"`
	LoggingLevel              *interface{}                                  `json:"loggingLevel,omitempty"`
	PackageConnectionManagers *map[string]map[string]SSISExecutionParameter `json:"packageConnectionManagers,omitempty"`
	PackageLocation           SSISPackageLocation                           `json:"packageLocation"`
	PackageParameters         *map[string]SSISExecutionParameter            `json:"packageParameters,omitempty"`
	ProjectConnectionManagers *map[string]map[string]SSISExecutionParameter `json:"projectConnectionManagers,omitempty"`
	ProjectParameters         *map[string]SSISExecutionParameter            `json:"projectParameters,omitempty"`
	PropertyOverrides         *map[string]SSISPropertyOverride              `json:"propertyOverrides,omitempty"`
	Runtime                   *interface{}                                  `json:"runtime,omitempty"`
}
