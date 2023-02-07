package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheConfiguration struct {
	CacheDuration                *string                    `json:"cacheDuration,omitempty"`
	DynamicCompression           *DynamicCompressionEnabled `json:"dynamicCompression,omitempty"`
	QueryParameterStripDirective *FrontDoorQuery            `json:"queryParameterStripDirective,omitempty"`
	QueryParameters              *string                    `json:"queryParameters,omitempty"`
}
