package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemRootSquashSettings struct {
	Mode             *AmlFilesystemSquashMode `json:"mode,omitempty"`
	NoSquashNidLists *string                  `json:"noSquashNidLists,omitempty"`
	SquashGID        *int64                   `json:"squashGID,omitempty"`
	SquashUID        *int64                   `json:"squashUID,omitempty"`
	Status           *string                  `json:"status,omitempty"`
}
