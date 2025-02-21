package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JupyterKernelConfig struct {
	Argv        *[]string `json:"argv,omitempty"`
	DisplayName *string   `json:"displayName,omitempty"`
	Language    *string   `json:"language,omitempty"`
}
