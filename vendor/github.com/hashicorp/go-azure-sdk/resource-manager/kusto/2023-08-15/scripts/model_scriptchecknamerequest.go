package scripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptCheckNameRequest struct {
	Name string     `json:"name"`
	Type ScriptType `json:"type"`
}
