package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemHsmSettings struct {
	Container        string  `json:"container"`
	ImportPrefix     *string `json:"importPrefix,omitempty"`
	LoggingContainer string  `json:"loggingContainer"`
}
