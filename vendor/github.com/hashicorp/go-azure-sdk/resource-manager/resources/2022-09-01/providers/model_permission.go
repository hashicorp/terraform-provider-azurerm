package providers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Permission struct {
	Actions        *[]string `json:"actions,omitempty"`
	DataActions    *[]string `json:"dataActions,omitempty"`
	NotActions     *[]string `json:"notActions,omitempty"`
	NotDataActions *[]string `json:"notDataActions,omitempty"`
}
