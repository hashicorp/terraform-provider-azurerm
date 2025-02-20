package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DistcpSettings struct {
	DistcpOptions           *string `json:"distcpOptions,omitempty"`
	ResourceManagerEndpoint string  `json:"resourceManagerEndpoint"`
	TempScriptPath          string  `json:"tempScriptPath"`
}
