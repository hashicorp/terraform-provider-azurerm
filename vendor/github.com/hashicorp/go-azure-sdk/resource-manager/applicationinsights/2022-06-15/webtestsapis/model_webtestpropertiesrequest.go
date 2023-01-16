package webtestsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestPropertiesRequest struct {
	FollowRedirects        *bool          `json:"FollowRedirects,omitempty"`
	HTTPVerb               *string        `json:"HttpVerb,omitempty"`
	Headers                *[]HeaderField `json:"Headers,omitempty"`
	ParseDependentRequests *bool          `json:"ParseDependentRequests,omitempty"`
	RequestBody            *string        `json:"RequestBody,omitempty"`
	RequestUrl             *string        `json:"RequestUrl,omitempty"`
}
