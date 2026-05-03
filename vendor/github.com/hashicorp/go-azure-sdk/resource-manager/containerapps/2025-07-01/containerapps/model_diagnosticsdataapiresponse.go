package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsDataApiResponse struct {
	RenderingProperties *DiagnosticRendering               `json:"renderingProperties,omitempty"`
	Table               *DiagnosticDataTableResponseObject `json:"table,omitempty"`
}
