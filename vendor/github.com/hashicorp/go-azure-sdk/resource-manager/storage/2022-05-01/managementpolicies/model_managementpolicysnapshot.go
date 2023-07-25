package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPolicySnapShot struct {
	Delete        *DateAfterCreation `json:"delete,omitempty"`
	TierToArchive *DateAfterCreation `json:"tierToArchive,omitempty"`
	TierToCool    *DateAfterCreation `json:"tierToCool,omitempty"`
}
