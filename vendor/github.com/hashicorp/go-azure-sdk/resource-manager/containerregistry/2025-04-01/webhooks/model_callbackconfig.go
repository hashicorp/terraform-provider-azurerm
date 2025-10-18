package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CallbackConfig struct {
	CustomHeaders *map[string]string `json:"customHeaders,omitempty"`
	ServiceUri    string             `json:"serviceUri"`
}
