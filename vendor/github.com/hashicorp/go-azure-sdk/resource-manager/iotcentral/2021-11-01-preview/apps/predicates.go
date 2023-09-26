package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p AppOperationPredicate) Matches(input App) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type AppTemplateOperationPredicate struct {
	Description     *string
	Industry        *string
	ManifestId      *string
	ManifestVersion *string
	Name            *string
	Order           *float64
	Title           *string
}

func (p AppTemplateOperationPredicate) Matches(input AppTemplate) bool {

	if p.Description != nil && (input.Description == nil || *p.Description != *input.Description) {
		return false
	}

	if p.Industry != nil && (input.Industry == nil || *p.Industry != *input.Industry) {
		return false
	}

	if p.ManifestId != nil && (input.ManifestId == nil || *p.ManifestId != *input.ManifestId) {
		return false
	}

	if p.ManifestVersion != nil && (input.ManifestVersion == nil || *p.ManifestVersion != *input.ManifestVersion) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Order != nil && (input.Order == nil || *p.Order != *input.Order) {
		return false
	}

	if p.Title != nil && (input.Title == nil || *p.Title != *input.Title) {
		return false
	}

	return true
}
