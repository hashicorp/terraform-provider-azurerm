package routes

type AfdRouteCacheConfiguration struct {
	CompressionSettings        *interface{}                   `json:"compressionSettings,omitempty"`
	QueryParameters            *string                        `json:"queryParameters,omitempty"`
	QueryStringCachingBehavior *AfdQueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
}
