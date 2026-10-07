package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthPlatform struct {
	ConfigFilePath *string `json:"configFilePath,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
	RuntimeVersion *string `json:"runtimeVersion,omitempty"`
}
