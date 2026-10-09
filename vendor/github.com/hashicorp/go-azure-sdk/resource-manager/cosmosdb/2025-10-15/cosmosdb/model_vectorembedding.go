package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorEmbedding struct {
	DataType         VectorDataType   `json:"dataType"`
	Dimensions       int64            `json:"dimensions"`
	DistanceFunction DistanceFunction `json:"distanceFunction"`
	Path             string           `json:"path"`
}
