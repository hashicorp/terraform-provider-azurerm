package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BigDataPoolParametrizationReference struct {
	ReferenceName interface{}              `json:"referenceName"`
	Type          BigDataPoolReferenceType `json:"type"`
}
