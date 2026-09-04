package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomEntityAlias struct {
	AccentSensitive   *bool  `json:"accentSensitive,omitempty"`
	CaseSensitive     *bool  `json:"caseSensitive,omitempty"`
	FuzzyEditDistance *int64 `json:"fuzzyEditDistance,omitempty"`
	Text              string `json:"text"`
}
