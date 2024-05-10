package managedhsmkeys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmAction struct {
	Type *KeyRotationPolicyActionType `json:"type,omitempty"`
}
