package dnsresolvers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubResourceListResult struct {
	NextLink *string        `json:"nextLink,omitempty"`
	Value    *[]SubResource `json:"value,omitempty"`
}
