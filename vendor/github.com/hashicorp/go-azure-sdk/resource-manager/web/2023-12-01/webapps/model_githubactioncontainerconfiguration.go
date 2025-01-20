package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitHubActionContainerConfiguration struct {
	ImageName *string `json:"imageName,omitempty"`
	Password  *string `json:"password,omitempty"`
	ServerURL *string `json:"serverUrl,omitempty"`
	Username  *string `json:"username,omitempty"`
}
