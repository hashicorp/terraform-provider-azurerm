package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotebookParameter struct {
	Type  *NotebookParameterType `json:"type,omitempty"`
	Value *interface{}           `json:"value,omitempty"`
}
