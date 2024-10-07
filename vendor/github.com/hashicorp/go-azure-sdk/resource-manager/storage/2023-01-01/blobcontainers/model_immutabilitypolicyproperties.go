package blobcontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutabilityPolicyProperties struct {
	Etag          *string                     `json:"etag,omitempty"`
	Properties    *ImmutabilityPolicyProperty `json:"properties,omitempty"`
	UpdateHistory *[]UpdateHistoryProperty    `json:"updateHistory,omitempty"`
}
