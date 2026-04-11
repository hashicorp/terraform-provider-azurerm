package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SingleVectorFieldResult struct {
	SearchScore      *float64 `json:"searchScore,omitempty"`
	VectorSimilarity *float64 `json:"vectorSimilarity,omitempty"`
}
