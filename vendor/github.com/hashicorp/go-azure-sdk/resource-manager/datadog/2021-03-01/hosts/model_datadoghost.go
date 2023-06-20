package hosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogHost struct {
	Aliases *[]string            `json:"aliases,omitempty"`
	Apps    *[]string            `json:"apps,omitempty"`
	Meta    *DatadogHostMetadata `json:"meta,omitempty"`
	Name    *string              `json:"name,omitempty"`
}
