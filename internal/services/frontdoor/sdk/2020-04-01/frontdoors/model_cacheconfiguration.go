package frontdoors

type CacheConfiguration struct {
	CacheDuration                *string                    `json:"cacheDuration,omitempty"`
	DynamicCompression           *DynamicCompressionEnabled `json:"dynamicCompression,omitempty"`
	QueryParameterStripDirective *FrontDoorQuery            `json:"queryParameterStripDirective,omitempty"`
	QueryParameters              *string                    `json:"queryParameters,omitempty"`
}
