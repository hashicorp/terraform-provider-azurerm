package nginxconfigurationanalysis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisDiagnostic struct {
	Description string  `json:"description"`
	Directive   string  `json:"directive"`
	File        string  `json:"file"`
	Id          *string `json:"id,omitempty"`
	Line        float64 `json:"line"`
	Message     string  `json:"message"`
	Rule        string  `json:"rule"`
}
