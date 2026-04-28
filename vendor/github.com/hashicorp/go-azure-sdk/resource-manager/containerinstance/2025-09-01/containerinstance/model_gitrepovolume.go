package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitRepoVolume struct {
	Directory  *string `json:"directory,omitempty"`
	Repository string  `json:"repository"`
	Revision   *string `json:"revision,omitempty"`
}
