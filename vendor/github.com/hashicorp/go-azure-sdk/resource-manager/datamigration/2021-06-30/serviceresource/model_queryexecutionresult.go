package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryExecutionResult struct {
	QueryText         *string              `json:"queryText,omitempty"`
	SourceResult      *ExecutionStatistics `json:"sourceResult,omitempty"`
	StatementsInBatch *int64               `json:"statementsInBatch,omitempty"`
	TargetResult      *ExecutionStatistics `json:"targetResult,omitempty"`
}
