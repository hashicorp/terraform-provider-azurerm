package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynapseNotebookReference struct {
	ReferenceName string                `json:"referenceName"`
	Type          NotebookReferenceType `json:"type"`
}
