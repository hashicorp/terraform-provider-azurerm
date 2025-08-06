package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActivityScriptBlock struct {
	Parameters *[]ScriptActivityParameter `json:"parameters,omitempty"`
	Text       interface{}                `json:"text"`
	Type       interface{}                `json:"type"`
}
