package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListNotebookKeysResult struct {
	PrimaryAccessKey   *string `json:"primaryAccessKey,omitempty"`
	SecondaryAccessKey *string `json:"secondaryAccessKey,omitempty"`
}
