package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptsToExecute struct {
	CreationScript *ScriptReference `json:"creationScript,omitempty"`
	StartupScript  *ScriptReference `json:"startupScript,omitempty"`
}
