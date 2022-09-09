package iotdpsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPFilterRule struct {
	Action     IPFilterActionType  `json:"action"`
	FilterName string              `json:"filterName"`
	IPMask     string              `json:"ipMask"`
	Target     *IPFilterTargetType `json:"target,omitempty"`
}
