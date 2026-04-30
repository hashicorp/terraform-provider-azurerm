package nginxconfigurationresponses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisResultData struct {
	Diagnostics *[]DiagnosticItem     `json:"diagnostics,omitempty"`
	Errors      *[]AnalysisDiagnostic `json:"errors,omitempty"`
}
