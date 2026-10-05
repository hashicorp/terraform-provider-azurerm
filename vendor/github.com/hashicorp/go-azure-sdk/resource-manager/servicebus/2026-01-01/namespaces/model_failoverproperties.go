package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailOverProperties struct {
	Force           *bool   `json:"force,omitempty"`
	PrimaryLocation *string `json:"primaryLocation,omitempty"`
}
