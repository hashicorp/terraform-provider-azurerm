package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticsProperties struct {
	DataProviderMetadata *DiagnosticDataProviderMetadata `json:"dataProviderMetadata,omitempty"`
	Dataset              *[]DiagnosticsDataApiResponse   `json:"dataset,omitempty"`
	Metadata             *DiagnosticsDefinition          `json:"metadata,omitempty"`
	Status               *DiagnosticsStatus              `json:"status,omitempty"`
}
