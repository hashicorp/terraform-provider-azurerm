package connectionanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteDiagnosticsSettings struct {
	Bandwidth           *float64 `json:"bandwidth,omitempty"`
	ErCircuitResourceId *string  `json:"erCircuitResourceId,omitempty"`
}
