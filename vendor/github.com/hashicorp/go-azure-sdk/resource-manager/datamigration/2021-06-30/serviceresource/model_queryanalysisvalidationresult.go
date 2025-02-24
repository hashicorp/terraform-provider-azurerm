package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryAnalysisValidationResult struct {
	QueryResults     *QueryExecutionResult `json:"queryResults,omitempty"`
	ValidationErrors *ValidationError      `json:"validationErrors,omitempty"`
}
