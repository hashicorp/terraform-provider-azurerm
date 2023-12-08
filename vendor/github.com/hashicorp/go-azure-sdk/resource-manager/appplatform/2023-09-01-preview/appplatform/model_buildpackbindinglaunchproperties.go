package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildpackBindingLaunchProperties struct {
	Properties *map[string]string `json:"properties,omitempty"`
	Secrets    *map[string]string `json:"secrets,omitempty"`
}
