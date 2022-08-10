package videoanalyzer

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Endpoint struct {
	EndpointUrl *string                   `json:"endpointUrl,omitempty"`
	Type        VideoAnalyzerEndpointType `json:"type"`
}
