package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppProbe struct {
	FailureThreshold              *int64                      `json:"failureThreshold,omitempty"`
	HTTPGet                       *ContainerAppProbeHTTPGet   `json:"httpGet,omitempty"`
	InitialDelaySeconds           *int64                      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds                 *int64                      `json:"periodSeconds,omitempty"`
	SuccessThreshold              *int64                      `json:"successThreshold,omitempty"`
	TcpSocket                     *ContainerAppProbeTcpSocket `json:"tcpSocket,omitempty"`
	TerminationGracePeriodSeconds *int64                      `json:"terminationGracePeriodSeconds,omitempty"`
	TimeoutSeconds                *int64                      `json:"timeoutSeconds,omitempty"`
	Type                          *Type                       `json:"type,omitempty"`
}
