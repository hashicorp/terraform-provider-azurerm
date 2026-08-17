package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionActiveSession struct {
	Protocol              *BastionConnectProtocol `json:"protocol,omitempty"`
	ResourceType          *string                 `json:"resourceType,omitempty"`
	SessionDurationInMins *float64                `json:"sessionDurationInMins,omitempty"`
	SessionId             *string                 `json:"sessionId,omitempty"`
	StartTime             *interface{}            `json:"startTime,omitempty"`
	TargetHostName        *string                 `json:"targetHostName,omitempty"`
	TargetIPAddress       *string                 `json:"targetIpAddress,omitempty"`
	TargetResourceGroup   *string                 `json:"targetResourceGroup,omitempty"`
	TargetResourceId      *string                 `json:"targetResourceId,omitempty"`
	TargetSubscriptionId  *string                 `json:"targetSubscriptionId,omitempty"`
	UserName              *string                 `json:"userName,omitempty"`
}
