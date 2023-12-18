package redis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCacheAccessPolicy struct {
	Id         *string                           `json:"id,omitempty"`
	Name       *string                           `json:"name,omitempty"`
	Properties *RedisCacheAccessPolicyProperties `json:"properties,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
