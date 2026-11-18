package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureActiveDirectoryLogin struct {
	DisableWWWAuthenticate *bool     `json:"disableWWWAuthenticate,omitempty"`
	LoginParameters        *[]string `json:"loginParameters,omitempty"`
}
