package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Volume struct {
	AzureFile *AzureFileVolume   `json:"azureFile,omitempty"`
	EmptyDir  *interface{}       `json:"emptyDir,omitempty"`
	GitRepo   *GitRepoVolume     `json:"gitRepo,omitempty"`
	Name      string             `json:"name"`
	Secret    *map[string]string `json:"secret,omitempty"`
}
