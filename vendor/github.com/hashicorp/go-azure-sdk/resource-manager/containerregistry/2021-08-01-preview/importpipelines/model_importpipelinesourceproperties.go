package importpipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportPipelineSourceProperties struct {
	KeyVaultUri string              `json:"keyVaultUri"`
	Type        *PipelineSourceType `json:"type,omitempty"`
	Uri         *string             `json:"uri,omitempty"`
}
