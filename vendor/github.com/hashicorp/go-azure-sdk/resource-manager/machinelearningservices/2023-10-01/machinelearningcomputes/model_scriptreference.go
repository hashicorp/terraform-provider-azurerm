package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptReference struct {
	ScriptArguments *string `json:"scriptArguments,omitempty"`
	ScriptData      *string `json:"scriptData,omitempty"`
	ScriptSource    *string `json:"scriptSource,omitempty"`
	Timeout         *string `json:"timeout,omitempty"`
}
