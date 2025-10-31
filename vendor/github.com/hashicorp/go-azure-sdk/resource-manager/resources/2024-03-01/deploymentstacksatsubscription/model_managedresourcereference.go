package deploymentstacksatsubscription

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedResourceReference struct {
	DenyStatus *DenyStatusMode     `json:"denyStatus,omitempty"`
	Id         *string             `json:"id,omitempty"`
	Status     *ResourceStatusMode `json:"status,omitempty"`
}
