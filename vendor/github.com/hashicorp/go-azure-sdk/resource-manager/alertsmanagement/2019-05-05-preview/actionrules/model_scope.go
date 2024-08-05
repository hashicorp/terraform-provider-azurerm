package actionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Scope struct {
	ScopeType *ScopeType `json:"scopeType,omitempty"`
	Values    *[]string  `json:"values,omitempty"`
}
