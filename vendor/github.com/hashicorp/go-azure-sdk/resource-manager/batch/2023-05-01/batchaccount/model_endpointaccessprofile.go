package batchaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointAccessProfile struct {
	DefaultAction EndpointAccessDefaultAction `json:"defaultAction"`
	IPRules       *[]IPRule                   `json:"ipRules,omitempty"`
}
