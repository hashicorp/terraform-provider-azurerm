package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectorInstructionModelBase struct {
	Parameters *interface{} `json:"parameters,omitempty"`
	Type       SettingType  `json:"type"`
}
