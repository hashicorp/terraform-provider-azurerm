package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitHubActionCodeConfiguration struct {
	RuntimeStack   *string `json:"runtimeStack,omitempty"`
	RuntimeVersion *string `json:"runtimeVersion,omitempty"`
}
