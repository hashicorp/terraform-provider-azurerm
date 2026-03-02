package dataflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerQueryTypeProperties struct {
	DocumentLocale *string             `json:"documentLocale,omitempty"`
	Script         *string             `json:"script,omitempty"`
	Sources        *[]PowerQuerySource `json:"sources,omitempty"`
}
