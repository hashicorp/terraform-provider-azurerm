package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InfrastructureProfile struct {
	BiosGuid                 *string       `json:"biosGuid,omitempty"`
	CheckpointType           *string       `json:"checkpointType,omitempty"`
	Checkpoints              *[]Checkpoint `json:"checkpoints,omitempty"`
	CloudId                  *string       `json:"cloudId,omitempty"`
	Generation               *int64        `json:"generation,omitempty"`
	InventoryItemId          *string       `json:"inventoryItemId,omitempty"`
	LastRestoredVMCheckpoint *Checkpoint   `json:"lastRestoredVMCheckpoint,omitempty"`
	TemplateId               *string       `json:"templateId,omitempty"`
	Uuid                     *string       `json:"uuid,omitempty"`
	VMmServerId              *string       `json:"vmmServerId,omitempty"`
	VirtualMachineName       *string       `json:"vmName,omitempty"`
}
