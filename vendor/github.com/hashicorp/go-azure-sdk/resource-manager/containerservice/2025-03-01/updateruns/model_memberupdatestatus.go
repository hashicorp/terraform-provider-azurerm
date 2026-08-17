package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MemberUpdateStatus struct {
	ClusterResourceId *string       `json:"clusterResourceId,omitempty"`
	Message           *string       `json:"message,omitempty"`
	Name              *string       `json:"name,omitempty"`
	OperationId       *string       `json:"operationId,omitempty"`
	Status            *UpdateStatus `json:"status,omitempty"`
}
