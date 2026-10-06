package application

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationProperties struct {
	AllowUpdates   *bool   `json:"allowUpdates,omitempty"`
	DefaultVersion *string `json:"defaultVersion,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
}
