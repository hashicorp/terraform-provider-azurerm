package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryRequestOptions struct {
	AllowPartialScopes       *bool                     `json:"allowPartialScopes,omitempty"`
	AuthorizationScopeFilter *AuthorizationScopeFilter `json:"authorizationScopeFilter,omitempty"`
	ResultFormat             *ResultFormat             `json:"resultFormat,omitempty"`
	Skip                     *int64                    `json:"$skip,omitempty"`
	SkipToken                *string                   `json:"$skipToken,omitempty"`
	Top                      *int64                    `json:"$top,omitempty"`
}
