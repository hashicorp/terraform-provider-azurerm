package skillsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomEntity struct {
	AccentSensitive          *bool                `json:"accentSensitive,omitempty"`
	Aliases                  *[]CustomEntityAlias `json:"aliases,omitempty"`
	CaseSensitive            *bool                `json:"caseSensitive,omitempty"`
	DefaultAccentSensitive   *bool                `json:"defaultAccentSensitive,omitempty"`
	DefaultCaseSensitive     *bool                `json:"defaultCaseSensitive,omitempty"`
	DefaultFuzzyEditDistance *int64               `json:"defaultFuzzyEditDistance,omitempty"`
	Description              *string              `json:"description,omitempty"`
	FuzzyEditDistance        *int64               `json:"fuzzyEditDistance,omitempty"`
	Id                       *string              `json:"id,omitempty"`
	Name                     string               `json:"name"`
	Subtype                  *string              `json:"subtype,omitempty"`
	Type                     *string              `json:"type,omitempty"`
}
