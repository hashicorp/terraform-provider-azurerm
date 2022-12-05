package transformations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformationProperties struct {
	Etag                *string  `json:"etag,omitempty"`
	Query               *string  `json:"query,omitempty"`
	StreamingUnits      *int64   `json:"streamingUnits,omitempty"`
	ValidStreamingUnits *[]int64 `json:"validStreamingUnits,omitempty"`
}
