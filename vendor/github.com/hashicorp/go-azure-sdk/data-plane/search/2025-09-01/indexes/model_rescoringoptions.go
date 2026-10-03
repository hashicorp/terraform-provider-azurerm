package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RescoringOptions struct {
	DefaultOversampling  *float64                                     `json:"defaultOversampling,omitempty"`
	EnableRescoring      *bool                                        `json:"enableRescoring,omitempty"`
	RescoreStorageMethod *VectorSearchCompressionRescoreStorageMethod `json:"rescoreStorageMethod,omitempty"`
}
