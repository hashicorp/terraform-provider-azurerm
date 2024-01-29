package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CSharpFunctionBindingProperties struct {
	Class      *string     `json:"class,omitempty"`
	DllPath    *string     `json:"dllPath,omitempty"`
	Method     *string     `json:"method,omitempty"`
	UpdateMode *UpdateMode `json:"updateMode,omitempty"`
}
