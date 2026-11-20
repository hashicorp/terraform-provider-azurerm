package edgeactionexecutionfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionExecutionFilterUpdateProperties struct {
	ExecutionFilterIdentifierHeaderName  *string `json:"executionFilterIdentifierHeaderName,omitempty"`
	ExecutionFilterIdentifierHeaderValue *string `json:"executionFilterIdentifierHeaderValue,omitempty"`
	VersionId                            *string `json:"versionId,omitempty"`
}
