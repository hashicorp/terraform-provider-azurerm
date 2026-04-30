package instancefailovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstancePairInfo struct {
	PartnerManagedInstanceId *string `json:"partnerManagedInstanceId,omitempty"`
	PrimaryManagedInstanceId *string `json:"primaryManagedInstanceId,omitempty"`
}
