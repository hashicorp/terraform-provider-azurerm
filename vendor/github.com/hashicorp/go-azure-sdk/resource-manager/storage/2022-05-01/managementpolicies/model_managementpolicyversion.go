package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPolicyVersion struct {
	Delete        *DateAfterCreation `json:"delete"`
	TierToArchive *DateAfterCreation `json:"tierToArchive"`
	TierToCool    *DateAfterCreation `json:"tierToCool"`
}
