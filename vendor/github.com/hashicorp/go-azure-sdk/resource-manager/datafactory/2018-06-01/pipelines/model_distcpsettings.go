package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DistcpSettings struct {
	DistcpOptions           *interface{} `json:"distcpOptions,omitempty"`
	ResourceManagerEndpoint interface{}  `json:"resourceManagerEndpoint"`
	TempScriptPath          interface{}  `json:"tempScriptPath"`
}
