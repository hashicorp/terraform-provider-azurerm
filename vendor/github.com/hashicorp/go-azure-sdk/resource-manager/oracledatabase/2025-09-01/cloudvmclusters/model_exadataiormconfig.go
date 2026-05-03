package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExadataIormConfig struct {
	DbPlans          *[]DbIormConfig     `json:"dbPlans,omitempty"`
	LifecycleDetails *string             `json:"lifecycleDetails,omitempty"`
	LifecycleState   *IormLifecycleState `json:"lifecycleState,omitempty"`
	Objective        *Objective          `json:"objective,omitempty"`
}
