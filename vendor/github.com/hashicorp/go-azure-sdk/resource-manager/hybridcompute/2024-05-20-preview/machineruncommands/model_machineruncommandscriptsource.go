package machineruncommands

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineRunCommandScriptSource struct {
	CommandId                *string                    `json:"commandId,omitempty"`
	Script                   *string                    `json:"script,omitempty"`
	ScriptUri                *string                    `json:"scriptUri,omitempty"`
	ScriptUriManagedIdentity *RunCommandManagedIdentity `json:"scriptUriManagedIdentity,omitempty"`
}
