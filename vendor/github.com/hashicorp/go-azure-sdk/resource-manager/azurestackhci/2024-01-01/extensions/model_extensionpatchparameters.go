package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionPatchParameters struct {
	EnableAutomaticUpgrade *bool        `json:"enableAutomaticUpgrade,omitempty"`
	ProtectedSettings      *interface{} `json:"protectedSettings,omitempty"`
	Settings               *interface{} `json:"settings,omitempty"`
	TypeHandlerVersion     *string      `json:"typeHandlerVersion,omitempty"`
}
