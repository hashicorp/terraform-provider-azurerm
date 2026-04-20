package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerQuerySinkMapping struct {
	DataflowSinks *[]PowerQuerySink `json:"dataflowSinks,omitempty"`
	QueryName     *string           `json:"queryName,omitempty"`
}
