package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionValueProperties struct {
	ExtensionType *string `json:"extensionType,omitempty"`
	Publisher     *string `json:"publisher,omitempty"`
	Version       *string `json:"version,omitempty"`
}
