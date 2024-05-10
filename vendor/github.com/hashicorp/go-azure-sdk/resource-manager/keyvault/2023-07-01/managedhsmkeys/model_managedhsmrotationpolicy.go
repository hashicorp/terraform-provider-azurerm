package managedhsmkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmRotationPolicy struct {
	Attributes      *ManagedHsmKeyRotationPolicyAttributes `json:"attributes,omitempty"`
	LifetimeActions *[]ManagedHsmLifetimeAction            `json:"lifetimeActions,omitempty"`
}
