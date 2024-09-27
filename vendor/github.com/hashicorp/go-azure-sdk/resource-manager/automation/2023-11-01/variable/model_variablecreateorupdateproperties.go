package variable

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VariableCreateOrUpdateProperties struct {
	Description *string `json:"description,omitempty"`
	IsEncrypted *bool   `json:"isEncrypted,omitempty"`
	Value       *string `json:"value,omitempty"`
}
