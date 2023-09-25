package virtualmachineruncommands

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunCommandDocumentBase struct {
	Description string               `json:"description"`
	Id          string               `json:"id"`
	Label       string               `json:"label"`
	OsType      OperatingSystemTypes `json:"osType"`
	Schema      string               `json:"$schema"`
}
