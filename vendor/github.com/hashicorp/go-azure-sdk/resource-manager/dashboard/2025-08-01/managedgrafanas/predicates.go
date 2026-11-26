package managedgrafanas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaAvailablePluginOperationPredicate struct {
	Author   *string
	Name     *string
	PluginId *string
	Type     *string
}

func (p GrafanaAvailablePluginOperationPredicate) Matches(input GrafanaAvailablePlugin) bool {

	if p.Author != nil && (input.Author == nil || *p.Author != *input.Author) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.PluginId != nil && (input.PluginId == nil || *p.PluginId != *input.PluginId) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type ManagedGrafanaOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ManagedGrafanaOperationPredicate) Matches(input ManagedGrafana) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
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
