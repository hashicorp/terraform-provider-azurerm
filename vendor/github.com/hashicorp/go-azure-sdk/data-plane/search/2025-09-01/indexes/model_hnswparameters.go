package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HnswParameters struct {
	EfConstruction *int64                       `json:"efConstruction,omitempty"`
	EfSearch       *int64                       `json:"efSearch,omitempty"`
	M              *int64                       `json:"m,omitempty"`
	Metric         *VectorSearchAlgorithmMetric `json:"metric,omitempty"`
}
