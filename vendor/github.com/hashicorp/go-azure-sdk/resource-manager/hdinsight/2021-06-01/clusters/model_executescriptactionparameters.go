package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecuteScriptActionParameters struct {
	PersistOnSuccess bool                   `json:"persistOnSuccess"`
	ScriptActions    *[]RuntimeScriptAction `json:"scriptActions,omitempty"`
}
