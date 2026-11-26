package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressConfiguration struct {
	HeaderCountLimit              *int64  `json:"headerCountLimit,omitempty"`
	RequestIdleTimeout            *int64  `json:"requestIdleTimeout,omitempty"`
	TerminationGracePeriodSeconds *int64  `json:"terminationGracePeriodSeconds,omitempty"`
	WorkloadProfileName           *string `json:"workloadProfileName,omitempty"`
}
