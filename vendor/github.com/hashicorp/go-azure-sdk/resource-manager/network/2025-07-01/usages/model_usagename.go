package usages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsageName struct {
	LocalizedValue *string `json:"localizedValue,omitempty"`
	Value          *string `json:"value,omitempty"`
}
