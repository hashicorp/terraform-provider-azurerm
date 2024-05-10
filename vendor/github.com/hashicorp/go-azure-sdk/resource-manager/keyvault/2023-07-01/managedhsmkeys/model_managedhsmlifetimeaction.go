package managedhsmkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmLifetimeAction struct {
	Action  *ManagedHsmAction  `json:"action,omitempty"`
	Trigger *ManagedHsmTrigger `json:"trigger,omitempty"`
}
