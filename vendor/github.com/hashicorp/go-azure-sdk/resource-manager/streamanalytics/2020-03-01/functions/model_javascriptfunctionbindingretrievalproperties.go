package functions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JavaScriptFunctionBindingRetrievalProperties struct {
	Script  *string  `json:"script,omitempty"`
	UdfType *UdfType `json:"udfType,omitempty"`
}
