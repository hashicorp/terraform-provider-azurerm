package vipswap

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwapResource struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *SwapResourceProperties `json:"properties,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
