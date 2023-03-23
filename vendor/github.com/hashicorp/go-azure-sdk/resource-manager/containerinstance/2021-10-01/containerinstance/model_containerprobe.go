package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerProbe struct {
	Exec                *ContainerExec    `json:"exec,omitempty"`
	FailureThreshold    *int64            `json:"failureThreshold,omitempty"`
	HTTPGet             *ContainerHTTPGet `json:"httpGet,omitempty"`
	InitialDelaySeconds *int64            `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int64            `json:"periodSeconds,omitempty"`
	SuccessThreshold    *int64            `json:"successThreshold,omitempty"`
	TimeoutSeconds      *int64            `json:"timeoutSeconds,omitempty"`
}
