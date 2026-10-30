package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandInput struct {
	CommandId  string                      `json:"commandId"`
	Parameters *[]RunCommandInputParameter `json:"parameters,omitempty"`
	Script     *[]string                   `json:"script,omitempty"`
}
