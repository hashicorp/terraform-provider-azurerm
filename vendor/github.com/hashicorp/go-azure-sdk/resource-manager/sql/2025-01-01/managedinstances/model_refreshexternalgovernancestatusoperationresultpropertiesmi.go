package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshExternalGovernanceStatusOperationResultPropertiesMI struct {
	ErrorMessage        *string `json:"errorMessage,omitempty"`
	ManagedInstanceName *string `json:"managedInstanceName,omitempty"`
	QueuedTime          *string `json:"queuedTime,omitempty"`
	RequestId           *string `json:"requestId,omitempty"`
	RequestType         *string `json:"requestType,omitempty"`
	Status              *string `json:"status,omitempty"`
}
