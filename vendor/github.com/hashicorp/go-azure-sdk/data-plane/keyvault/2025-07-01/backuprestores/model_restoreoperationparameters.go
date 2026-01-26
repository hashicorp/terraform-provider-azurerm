package backuprestores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreOperationParameters struct {
	FolderToRestore    string            `json:"folderToRestore"`
	SasTokenParameters SASTokenParameter `json:"sasTokenParameters"`
}
