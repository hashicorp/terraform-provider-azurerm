package exportpipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportPipelineTargetProperties struct {
	KeyVaultUri string  `json:"keyVaultUri"`
	Type        *string `json:"type,omitempty"`
	Uri         *string `json:"uri,omitempty"`
}
