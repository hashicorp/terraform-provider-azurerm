package rules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheConfiguration struct {
	CacheBehavior              *RuleCacheBehavior              `json:"cacheBehavior,omitempty"`
	CacheDuration              *string                         `json:"cacheDuration,omitempty"`
	IsCompressionEnabled       *RuleIsCompressionEnabled       `json:"isCompressionEnabled,omitempty"`
	QueryParameters            *string                         `json:"queryParameters,omitempty"`
	QueryStringCachingBehavior *RuleQueryStringCachingBehavior `json:"queryStringCachingBehavior,omitempty"`
}
