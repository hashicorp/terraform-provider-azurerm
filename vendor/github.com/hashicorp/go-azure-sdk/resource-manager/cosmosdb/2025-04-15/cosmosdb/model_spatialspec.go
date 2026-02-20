package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpatialSpec struct {
	Path  *string        `json:"path,omitempty"`
	Types *[]SpatialType `json:"types,omitempty"`
}
